/*
 * linapro-live-manage viewer H5 player.
 * Self-contained ES5-compatible script: parses room/tenant query params,
 * fetches the public play info once, renders one of four states
 * (live / replay / preview / none), attaches hls.js or native HLS playback,
 * and polls every 30 seconds while the live has not started or nothing is
 * found. Text is Chinese-only by documented decision. The API endpoint is
 * derived from the /h5 prefix of the current path so reverse-proxy sub-path
 * deployments keep working.
 */
(function () {
  'use strict';

  var POLL_INTERVAL_MS = 30000;
  var FETCH_TIMEOUT_MS = 10000;
  var HLS_RECOVERY_LIMIT = 2;

  var dom = {
    loadingSection: document.getElementById('loading-section'),
    playerSection: document.getElementById('player-section'),
    video: document.getElementById('video'),
    liveBadge: document.getElementById('live-badge'),
    replayBadge: document.getElementById('replay-badge'),
    playerError: document.getElementById('player-error'),
    playerErrorText: document.getElementById('player-error-text'),
    playerRetry: document.getElementById('player-retry'),
    unmuteHint: document.getElementById('unmute-hint'),
    coverSection: document.getElementById('cover-section'),
    coverImage: document.getElementById('cover-image'),
    coverFallback: document.getElementById('cover-fallback'),
    coverDate: document.getElementById('cover-date'),
    infoSection: document.getElementById('info-section'),
    stateChip: document.getElementById('state-chip'),
    liveRoom: document.getElementById('live-room'),
    liveTitle: document.getElementById('live-title'),
    liveMeta: document.getElementById('live-meta'),
    program: document.getElementById('program'),
    noticeSection: document.getElementById('notice-section'),
    noticeIcon: document.getElementById('notice-icon'),
    noticeText: document.getElementById('notice-text'),
    retryButton: document.getElementById('retry-button')
  };

  var state = { pollTimer: null, polling: false, hls: null, lastUrl: '', retryCount: 0 };

  /* ---- 参数与地址 ---- */

  function parseQuery() {
    var params = new URLSearchParams(window.location.search);
    return {
      room: (params.get('room') || '').trim(),
      tenant: (params.get('tenant') || '').trim()
    };
  }

  function apiEndpoint() {
    var path = window.location.pathname;
    var marker = path.lastIndexOf('/h5');
    if (marker >= 0) {
      path = path.slice(0, marker);
    } else {
      path = path.replace(/\/[^/]*$/, '');
    }
    if (path.charAt(path.length - 1) === '/') {
      path = path.slice(0, -1);
    }
    return path + '/api/v1/play';
  }

  /* ---- 数据加载 ---- */

  /* fetch with a hard timeout so a stalled connection cannot hang the page. */
  function requestWithTimeout(url) {
    var controller = typeof window.AbortController === 'function' ? new window.AbortController() : null;
    var init = { method: 'GET', cache: 'no-store' };
    var timer = null;
    if (controller) {
      init.signal = controller.signal;
      timer = window.setTimeout(function () {
        controller.abort();
      }, FETCH_TIMEOUT_MS);
    }
    return fetch(url, init).then(
      function (response) {
        if (timer) {
          window.clearTimeout(timer);
        }
        return response;
      },
      function (error) {
        if (timer) {
          window.clearTimeout(timer);
        }
        throw error;
      }
    );
  }

  function fetchPlayInfo() {
    var query = parseQuery();
    if (!query.room) {
      return Promise.reject({ kind: 'params' });
    }
    var url = apiEndpoint() +
      '?roomCode=' + encodeURIComponent(query.room) +
      (query.tenant !== '' ? '&tenantId=' + encodeURIComponent(query.tenant) : '');
    return requestWithTimeout(url)
      .then(function (response) {
        return response.json().then(function (payload) {
          return { status: response.status, payload: payload };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.status >= 200 && result.status < 300 && payload.code === 0 && payload.data) {
          return { kind: 'ok', play: payload.data };
        }
        var errorCode = payload.errorCode || '';
        if (errorCode === 'LIVE_MANAGE_PLAY_NOT_FOUND') {
          return { kind: 'none' };
        }
        if (errorCode === 'LIVE_MANAGE_PLAY_TENANT_REQUIRED' || errorCode === 'LIVE_MANAGE_PLAY_TENANT_INVALID') {
          return { kind: 'tenant', message: payload.message };
        }
        return { kind: 'error', message: payload.message || ('加载失败 (' + result.status + ')') };
      })
      .catch(function () {
        return { kind: 'network' };
      });
  }

  /* ---- 渲染 ---- */

  function hideAll() {
    dom.loadingSection.classList.add('hidden');
    dom.playerSection.classList.add('hidden');
    dom.coverSection.classList.add('hidden');
    dom.infoSection.classList.add('hidden');
    dom.noticeSection.classList.add('hidden');
    dom.liveBadge.classList.add('hidden');
    dom.replayBadge.classList.add('hidden');
    dom.playerError.classList.add('hidden');
    dom.unmuteHint.classList.add('hidden');
    dom.liveRoom.classList.add('hidden');
    dom.coverDate.classList.add('hidden');
    dom.liveTitle.textContent = '';
    dom.liveMeta.innerHTML = '';
    dom.program.innerHTML = '';
  }

  function stopPolling() {
    if (state.pollTimer) {
      window.clearInterval(state.pollTimer);
      state.pollTimer = null;
    }
    state.polling = false;
  }

  function startPolling() {
    if (state.pollTimer) {
      return;
    }
    state.polling = true;
    state.pollTimer = window.setInterval(function () {
      refresh();
    }, POLL_INTERVAL_MS);
  }

  function destroyPlayer() {
    clearRecoveryWatch();
    if (state.hls) {
      state.hls.destroy();
      state.hls = null;
    }
    dom.video.removeAttribute('src');
    try {
      dom.video.load();
    } catch (e) { /* 忽略空源加载异常 */ }
    dom.unmuteHint.classList.add('hidden');
  }

  function showPlayer(play) {
    dom.coverSection.classList.add('hidden');
    dom.noticeSection.classList.add('hidden');
    dom.playerSection.classList.remove('hidden');
    dom.infoSection.classList.remove('hidden');
    if (play.state === 2) {
      dom.replayBadge.classList.remove('hidden');
      setStateChip('replay', '已结束 · 回放');
    } else {
      dom.liveBadge.classList.remove('hidden');
      setStateChip('live', '直播中');
    }
    renderHeader(play);
    renderProgram(play);
    attachStream(play.liveUrl);
  }

  function showCover(play) {
    dom.playerSection.classList.add('hidden');
    dom.noticeSection.classList.add('hidden');
    dom.playerError.classList.add('hidden');
    dom.coverSection.classList.remove('hidden');
    dom.infoSection.classList.remove('hidden');
    setStateChip('preview', '未开始 · 预告');
    renderHeader(play);
    renderProgram(play);
    if (play.coverUrl) {
      dom.coverImage.src = play.coverUrl;
      dom.coverImage.style.display = 'block';
      dom.coverFallback.style.display = 'none';
    } else {
      dom.coverImage.style.display = 'none';
      dom.coverFallback.style.display = 'flex';
    }
    if (play.liveDate) {
      dom.coverDate.textContent = play.liveDate + ' 开播';
      dom.coverDate.classList.remove('hidden');
    } else {
      dom.coverDate.classList.add('hidden');
    }
  }

  /* ---- 提示区 SVG 图标（大厂风格，避免 emoji 渲染差异） ---- */

  var NOTICE_ICONS = {
    lock: '<svg viewBox="0 0 24 24" fill="none" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="4.5" y="10.5" width="15" height="10" rx="2.5"/><path d="M8 10.5V7.8a4 4 0 0 1 8 0v2.7"/><circle cx="12" cy="15.4" r="1.5"/></svg>',
    link: '<svg viewBox="0 0 24 24" fill="none" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M10 14a4.6 4.6 0 0 0 6.6.3l2.5-2.5a4.6 4.6 0 0 0-6.5-6.5l-1.4 1.4"/><path d="M14 10a4.6 4.6 0 0 0-6.6-.3l-2.5 2.5a4.6 4.6 0 0 0 6.5 6.5l1.4-1.4"/></svg>',
    signal: '<svg viewBox="0 0 24 24" fill="none" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 10.5a11.5 11.5 0 0 1 16 0"/><path d="M7 13.8a7.2 7.2 0 0 1 10 0"/><circle cx="12" cy="18.2" r="1.4" fill="currentColor" stroke="none"/></svg>'
  };

  function setNoticeIcon(name) {
    dom.noticeIcon.innerHTML = NOTICE_ICONS[name] || NOTICE_ICONS.signal;
  }

  function showNotice(kind, message) {
    stopPolling();
    destroyPlayer();
    dom.playerSection.classList.add('hidden');
    dom.infoSection.classList.add('hidden');
    dom.playerError.classList.add('hidden');
    dom.coverDate.classList.add('hidden');
    dom.coverSection.classList.remove('hidden');
    dom.coverImage.style.display = 'none';
    dom.coverFallback.style.display = 'flex';
    dom.noticeSection.classList.remove('hidden');
    if (kind === 'tenant') {
      setNoticeIcon('lock');
      dom.noticeText.textContent = message || '直播链接无效或已过期';
      dom.retryButton.classList.add('hidden');
      return;
    }
    if (kind === 'params') {
      setNoticeIcon('link');
      dom.noticeText.textContent = '缺少直播间参数，请通过分享链接打开';
      dom.retryButton.classList.add('hidden');
      return;
    }
    if (kind === 'network') {
      setNoticeIcon('signal');
      dom.noticeText.textContent = '网络异常，请稍后重试';
      dom.retryButton.classList.remove('hidden');
      startPolling();
      return;
    }
    setNoticeIcon('signal');
    dom.noticeText.textContent = '暂无直播，敬请期待';
    dom.retryButton.classList.remove('hidden');
    startPolling();
  }

  function setStateChip(theme, text) {
    dom.stateChip.className = 'state-chip ' + theme;
    dom.stateChip.textContent = text;
  }

  function renderHeader(play) {
    if (play.roomName) {
      dom.liveRoom.textContent = play.roomName;
      dom.liveRoom.classList.remove('hidden');
    } else {
      dom.liveRoom.textContent = '';
      dom.liveRoom.classList.add('hidden');
    }
    dom.liveTitle.textContent = play.title || play.roomName || '直播';
    var parts = [];
    if (play.liveDate) {
      parts.push(escapeHtml(play.liveDate));
    }
    if (play.startTime) {
      parts.push(escapeHtml(formatTime(play.startTime)));
    }
    if (parts.length) {
      dom.liveMeta.innerHTML = parts.join('<span class="sep"></span>');
    } else {
      dom.liveMeta.innerHTML = '';
    }
  }

  function renderProgram(play) {
    var cards = [];
    var songs = parseSongList(play.songList);
    var rows = [];

    if (songs.length) {
      var items = songs.map(function (song) {
        return '<li>' +
          '<span class="song-order">' + escapeHtml(String(song.order || '')) + '</span>' +
          '<span>' + escapeHtml(song.name || '') + '</span>' +
          (song.singer ? '<span class="song-singer">' + escapeHtml(song.singer) + '</span>' : '') +
          '</li>';
      }).join('');
      cards.push('<div class="program-card"><h3>诗歌节目单</h3><ul class="program-list">' + items + '</ul></div>');
    }

    if (play.songName) {
      rows.push(programRow('歌曲', play.songName + (play.leadSinger ? '（领唱：' + play.leadSinger + '）' : '')));
    } else if (play.leadSinger) {
      rows.push(programRow('领唱', play.leadSinger));
    }
    if (play.host) {
      rows.push(programRow('主持', play.host));
    }
    if (play.accompaniment) {
      rows.push(programRow('伴奏', play.accompaniment));
    }
    if (play.sermonTitle) {
      rows.push(programRow('讲道', play.sermonTitle + (play.preacher ? ' · ' + play.preacher : '')));
    } else if (play.preacher) {
      rows.push(programRow('讲员', play.preacher + (play.preacherIdentity ? '（' + play.preacherIdentity + '）' : '')));
    }
    if (play.scriptureRef) {
      rows.push(programRow('经文', play.scriptureRef));
    }
    if (play.scriptureContent) {
      rows.push(programRow('经文内容', play.scriptureContent));
    }
    if (play.outline) {
      rows.push(programRow('大纲', play.outline));
    }
    if (rows.length) {
      cards.push('<div class="program-card"><h3>聚会安排</h3>' + rows.join('') + '</div>');
    }
    dom.program.innerHTML = cards.join('');
  }

  function programRow(key, value) {
    return '<div class="program-row"><span class="k">' + escapeHtml(key) + '</span><span class="v">' + escapeHtml(value) + '</span></div>';
  }

  function parseSongList(raw) {
    if (!raw) {
      return [];
    }
    try {
      var parsed = JSON.parse(raw);
      if (!Array.isArray(parsed)) {
        return [];
      }
      return parsed.filter(function (song) {
        return song && (song.name || song.singer);
      });
    } catch (e) {
      return [];
    }
  }

  function formatTime(milli) {
    var date = new Date(Number(milli));
    if (isNaN(date.getTime())) {
      return '';
    }
    var hour = ('0' + date.getHours()).slice(-2);
    var minute = ('0' + date.getMinutes()).slice(-2);
    return hour + ':' + minute;
  }

  function escapeHtml(value) {
    return String(value == null ? '' : value)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  /* ---- 播放器 ---- */

  function attachStream(url) {
    if (!url) {
      showPlayerError('该直播暂未提供播放地址');
      return;
    }
    state.lastUrl = url;
    state.retryCount = 0;
    dom.video.muted = false;
    var video = dom.video;
    if (window.Hls && window.Hls.isSupported()) {
      var hls = new window.Hls({ maxBufferLength: 30 });
      hls.loadSource(url);
      hls.attachMedia(video);
      hls.on(window.Hls.Events.MANIFEST_PARSED, function () {
        tryPlay(video);
      });
      /* 分级恢复：网络类错误先续流（带超时监视），媒体类错误软恢复，
         超出次数才降级为可见的重试入口（头部播放器标准策略）。 */
      hls.on(window.Hls.Events.ERROR, function (_event, data) {
        if (!data || !data.fatal) {
          return;
        }
        if (state.retryCount >= HLS_RECOVERY_LIMIT) {
          showPlayerError('播放异常，请稍后重试');
          return;
        }
        state.retryCount += 1;
        if (data.type === window.Hls.ErrorTypes.NETWORK_ERROR) {
          hls.startLoad();
          /* startLoad 对"源彻底不可达"可能静默失败：10s 内无任何分片
             到达则判定恢复失败，降级为重试入口。 */
          armRecoveryWatch(video);
        } else if (data.type === window.Hls.ErrorTypes.MEDIA_ERROR) {
          hls.recoverMediaError();
        } else {
          showPlayerError('播放异常，请稍后重试');
        }
      });
      hls.on(window.Hls.Events.FRAG_LOADED, function () {
        /* 有数据到达即视为恢复成功 */
        clearRecoveryWatch();
      });
      state.hls = hls;
      return;
    }
    if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = url;
      tryPlay(video);
      return;
    }
    showPlayerError('当前设备或浏览器不支持 HLS 播放');
  }

  var recoveryWatchTimer = null;

  /* startLoad 的恢复失败盲区监视：源彻底不可达时 hls.js 不再抛 fatal，
     只能靠超时观察 video 是否有实际进展。 */
  function armRecoveryWatch(video) {
    clearRecoveryWatch();
    recoveryWatchTimer = window.setTimeout(function () {
      recoveryWatchTimer = null;
      if (video.readyState > 0 || !state.hls) {
        return;
      }
      showPlayerError('播放异常，请稍后重试');
    }, 10000);
  }

  function clearRecoveryWatch() {
    if (recoveryWatchTimer) {
      window.clearTimeout(recoveryWatchTimer);
      recoveryWatchTimer = null;
    }
  }

  /* 自动播放被浏览器拦截时降级为静音自动播放，并提示用户一键开启声音。 */
  function tryPlay(video) {
    video.play().then(function () {
      dom.unmuteHint.classList.add('hidden');
    }).catch(function () {
      video.muted = true;
      video.play().then(function () {
        dom.unmuteHint.classList.remove('hidden');
      }).catch(function () {
        /* 静音后仍失败：保持手动控制 */
      });
    });
  }

  function showPlayerError(message) {
    stopPolling();
    dom.playerErrorText.textContent = message;
    dom.playerError.classList.remove('hidden');
  }

  /* ---- 主流程 ---- */

  function refresh() {
    return fetchPlayInfo().then(function (result) {
      dom.loadingSection.classList.add('hidden');
      destroyPlayer();
      dom.playerError.classList.add('hidden');
      stopPolling();
      switch (result.kind) {
        case 'ok':
          if (result.play.state === 0) {
            showCover(result.play);
            startPolling();
          } else {
            showPlayer(result.play);
          }
          break;
        case 'none':
          showNotice('none');
          break;
        case 'tenant':
          showNotice('tenant', result.message);
          break;
        case 'params':
          showNotice('params');
          break;
        default:
          showNotice(result.kind, result.message);
      }
    });
  }

  /* ---- 事件 ---- */

  dom.retryButton.addEventListener('click', refresh);

  dom.playerRetry.addEventListener('click', function () {
    dom.playerError.classList.add('hidden');
    attachStream(state.lastUrl);
  });

  dom.unmuteHint.addEventListener('click', function () {
    dom.video.muted = false;
    dom.unmuteHint.classList.add('hidden');
    dom.video.play().catch(function () { /* 用户已手动控制 */ });
  });

  /* 页面从后台回到前台时，轮询态（预告/无直播/网络异常）立即刷新，
     播放态不打断，交由 hls.js 自动回到直播边缘。 */
  document.addEventListener('visibilitychange', function () {
    if (document.visibilityState === 'visible' && state.polling) {
      refresh();
    }
  });

  refresh();
})();
