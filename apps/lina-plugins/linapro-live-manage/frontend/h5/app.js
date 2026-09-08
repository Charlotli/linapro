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
  var HEARTBEAT_INTERVAL_MS = 30000;
  var SESSION_STORAGE_KEY = 'linapro-live-view-session';
  var REPLAY_PAGE_SIZE = 50;

  var dom = {
    loadingSection: document.getElementById('loading-section'),
    playerSection: document.getElementById('player-section'),
    video: document.getElementById('video'),
    liveBadge: document.getElementById('live-badge'),
    replayBadge: document.getElementById('replay-badge'),
    viewersBadge: document.getElementById('viewers-badge'),
    viewersCount: document.getElementById('viewers-count'),
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
    calendarBar: document.getElementById('calendar-bar'),
    calendarOpen: document.getElementById('calendar-open'),
    calendarCopy: document.getElementById('calendar-copy'),
    shareBar: document.getElementById('share-bar'),
    shareButton: document.getElementById('share-button'),
    shareButtonLabel: document.getElementById('share-button-label'),
    replaySection: document.getElementById('replay-section'),
    replayList: document.getElementById('replay-list'),
    toolRow: document.getElementById('tool-row'),
    toolAnnouncement: document.getElementById('tool-announcement'),
    toolBible: document.getElementById('tool-bible'),
    announcementOverlay: document.getElementById('announcement-overlay'),
    announcementList: document.getElementById('announcement-list'),
    bibleOverlay: document.getElementById('bible-overlay'),
    bibleBody: document.getElementById('bible-body'),
    bibleTitle: document.getElementById('bible-title'),
    bibleBack: document.getElementById('bible-back'),
    bibleNav: document.getElementById('bible-nav'),
    biblePrev: document.getElementById('bible-prev'),
    bibleNext: document.getElementById('bible-next'),
    toast: document.getElementById('toast'),
    noticeSection: document.getElementById('notice-section'),
    noticeIcon: document.getElementById('notice-icon'),
    noticeText: document.getElementById('notice-text'),
    retryButton: document.getElementById('retry-button')
  };

  var state = {
    pollTimer: null,
    polling: false,
    hls: null,
    lastUrl: '',
    retryCount: 0,
    toastTimer: null,
    heartbeatTimer: null,
    replayItems: [],
    presentable: false,
    announcements: [],
    bibleBooks: null,
    bibleBookSn: null,
    bibleChapter: null,
    bibleView: 'books'
  };

  /* ---- 参数与地址 ---- */

  function parseQuery() {
    var params = new URLSearchParams(window.location.search);
    return {
      room: (params.get('room') || '').trim(),
      tenant: (params.get('tenant') || '').trim()
    };
  }

  function apiBase() {
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
    return path;
  }

  function apiEndpoint() {
    return apiBase() + '/api/v1/play';
  }

  function subscribeURL() {
    var query = parseQuery();
    return apiBase() + '/api/v1/subscribe' +
      '?roomCode=' + encodeURIComponent(query.room) +
      (query.tenant !== '' ? '&tenantId=' + encodeURIComponent(query.tenant) : '');
  }

  function replaysEndpoint() {
    var query = parseQuery();
    return apiBase() + '/api/v1/replays' +
      '?roomCode=' + encodeURIComponent(query.room) +
      '&page=1&pageSize=' + REPLAY_PAGE_SIZE +
      (query.tenant !== '' ? '&tenantId=' + encodeURIComponent(query.tenant) : '');
  }

  function heartbeatEndpoint() {
    return apiBase() + '/api/v1/view/heartbeat';
  }

  function announcementsEndpoint() {
    var query = parseQuery();
    return apiBase() + '/api/v1/announcements' +
      '?roomCode=' + encodeURIComponent(query.room) +
      (query.tenant !== '' ? '&tenantId=' + encodeURIComponent(query.tenant) : '');
  }

  function bibleBooksEndpoint() {
    return apiBase() + '/api/v1/bible/books';
  }

  function bibleChapterEndpoint(volumeSn, chapter) {
    return apiBase() + '/api/v1/bible/chapter' +
      '?volumeSn=' + encodeURIComponent(volumeSn) +
      '&chapter=' + encodeURIComponent(chapter);
  }

  /* ---- 公告面板 ---- */

  /* 拉取当前直播间的启用公告；网络或业务失败返回 null（与真实空列表
     区分，面板据此显示可重试的失败态而非误导性的"暂无公告"）。 */
  function fetchAnnouncements() {
    var query = parseQuery();
    if (!query.room) {
      return Promise.resolve([]);
    }
    return requestWithTimeout(announcementsEndpoint())
      .then(function (response) {
        return response.json().then(function (payload) {
          return { status: response.status, payload: payload };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.status >= 200 && result.status < 300 && payload.code === 0 && payload.data) {
          return payload.data.list || [];
        }
        return null;
      })
      .catch(function () {
        return null;
      });
  }

  function renderAnnouncements(failed) {
    if (failed) {
      dom.announcementList.innerHTML = '<button type="button" class="sheet-retry">公告加载失败，轻触重试</button>';
      return;
    }
    if (!state.announcements.length) {
      dom.announcementList.innerHTML = '<div class="sheet-empty">暂无公告</div>';
      return;
    }
    dom.announcementList.innerHTML = state.announcements.map(function (item) {
      var timeHtml = '';
      var milli = Number(item.updatedAt);
      if (milli > 0) {
        var d = new Date(milli);
        var pad = function (n) { return ('0' + n).slice(-2); };
        timeHtml = '<div class="announcement-time">' +
          d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()) +
          ' 更新</div>';
      }
      return '<article class="announcement-card">' +
        '<h4 class="announcement-title">' + escapeHtml(item.title || '公告') + '</h4>' +
        '<div class="announcement-content">' + renderMultiline(item.content) + '</div>' +
        timeHtml +
        '</article>';
    }).join('');
  }

  function openAnnouncementSheet() {
    dom.announcementOverlay.classList.remove('hidden');
    dom.announcementList.innerHTML = '<div class="sheet-empty">加载中…</div>';
    dom.announcementList.scrollTop = 0;
    fetchAnnouncements().then(function (items) {
      state.announcements = items || [];
      renderAnnouncements(items === null);
    });
  }

  /* ---- 圣经面板：书卷 → 章节 → 阅读 ---- */

  /* 书卷目录整页缓存；章节按 (卷,章) 缓存上次内容，重复翻页不重复请求。 */
  function fetchBibleBooks() {
    if (state.bibleBooks) {
      return Promise.resolve(state.bibleBooks);
    }
    return requestWithTimeout(bibleBooksEndpoint())
      .then(function (response) {
        return response.json().then(function (payload) {
          return { status: response.status, payload: payload };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.status >= 200 && result.status < 300 && payload.code === 0 && payload.data) {
          state.bibleBooks = payload.data.list || [];
          return state.bibleBooks;
        }
        return null;
      })
      .catch(function () {
        return null;
      });
  }

  function chapterCacheKey(volumeSn, chapter) {
    return volumeSn + ':' + chapter;
  }

  function fetchBibleChapter(volumeSn, chapter) {
    var cached = state.bibleChapters && state.bibleChapters[chapterCacheKey(volumeSn, chapter)];
    if (cached) {
      return Promise.resolve(cached);
    }
    return requestWithTimeout(bibleChapterEndpoint(volumeSn, chapter))
      .then(function (response) {
        return response.json().then(function (payload) {
          return { status: response.status, payload: payload };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.status >= 200 && result.status < 300 && payload.code === 0 && payload.data) {
          var data = payload.data;
          if (data && data.list && data.list.length) {
            if (!state.bibleChapters) {
              state.bibleChapters = {};
            }
            state.bibleChapters[chapterCacheKey(volumeSn, chapter)] = data;
          }
          return data;
        }
        return { list: [], book: '', chapter: chapter };
      })
      .catch(function () {
        return { list: [], book: '', chapter: chapter, failed: true };
      });
  }

  function bookBySn(sn) {
    var books = state.bibleBooks || [];
    for (var i = 0; i < books.length; i++) {
      if (books[i].sn === sn) {
        return books[i];
      }
    }
    return null;
  }

  function renderBible() {
    if (state.bibleView === 'books') {
      dom.bibleTitle.textContent = '圣经';
      dom.bibleBack.classList.add('hidden');
      dom.bibleNav.classList.add('hidden');
      renderBibleBooks();
      dom.bibleBody.scrollTop = 0;
      return;
    }
    if (state.bibleView === 'chapters') {
      var book = bookBySn(state.bibleBookSn);
      dom.bibleTitle.textContent = book ? book.fullName : '圣经';
      dom.bibleBack.classList.remove('hidden');
      dom.bibleNav.classList.add('hidden');
      renderBibleChapters(book);
      dom.bibleBody.scrollTop = 0;
      return;
    }
    var reading = bookBySn(state.bibleBookSn);
    dom.bibleTitle.textContent = reading ? reading.fullName + ' 第' + state.bibleChapter + '章' : '圣经';
    dom.bibleBack.classList.remove('hidden');
    dom.bibleNav.classList.remove('hidden');
    dom.bibleBody.innerHTML = '<div class="sheet-empty">加载中…</div>';
    dom.bibleBody.scrollTop = 0;
    fetchBibleChapter(state.bibleBookSn, state.bibleChapter).then(function (data) {
      renderBibleVerses(data);
    });
  }

  function renderBibleBooks() {
    var books = state.bibleBooks;
    if (books === null || books === undefined) {
      dom.bibleBody.innerHTML = '<button type="button" class="sheet-retry">圣经加载失败，轻触重试</button>';
      return;
    }
    if (!books.length) {
      dom.bibleBody.innerHTML = '<div class="sheet-empty">暂无书卷数据</div>';
      return;
    }
    var sections = [];
    [[1, '旧约'], [2, '新约']].forEach(function (group) {
      var chips = books.filter(function (book) {
        return book.testament === group[0];
      }).map(function (book) {
        return '<button type="button" class="bible-book" data-sn="' + book.sn + '" title="' + escapeHtml(book.fullName) + '">' +
          escapeHtml(book.shortName) + '</button>';
      }).join('');
      sections.push('<div class="bible-group"><h4 class="bible-group-title">' + group[1] + '</h4>' +
        '<div class="bible-grid">' + chips + '</div></div>');
    });
    dom.bibleBody.innerHTML = sections.join('');
  }

  function renderBibleChapters(book) {
    if (!book) {
      dom.bibleBody.innerHTML = '<div class="sheet-empty">请选择书卷</div>';
      return;
    }
    var chips = [];
    for (var i = 1; i <= book.chapterCount; i++) {
      var current = state.bibleChapter === i ? ' data-current="1"' : '';
      chips.push('<button type="button" class="bible-chapter"' + current + ' data-chapter="' + i + '">' + i + '</button>');
    }
    dom.bibleBody.innerHTML = '<div class="bible-grid">' + chips.join('') + '</div>';
  }

  function renderBibleVerses(data) {
    var verses = (data && data.list) || [];
    if (data && data.failed) {
      dom.bibleBody.innerHTML = '<button type="button" class="sheet-retry">本章加载失败，轻触重试</button>';
      return;
    }
    if (!verses.length) {
      dom.bibleBody.innerHTML = '<div class="sheet-empty">本章暂无经文</div>';
      return;
    }
    dom.bibleBody.innerHTML = '<div class="bible-verses">' + verses.map(function (verse) {
      return '<p class="bible-verse"><span class="bible-verse-no">' + verse.verseSn + '</span>' + escapeHtml(verse.lection || '') + '</p>';
    }).join('') + '</div>';
  }

  function openBibleSheet() {
    dom.bibleOverlay.classList.remove('hidden');
    dom.bibleBody.innerHTML = '<div class="sheet-empty">加载中…</div>';
    dom.bibleBody.scrollTop = 0;
    state.bibleView = 'books';
    state.bibleBookSn = null;
    state.bibleChapter = null;
    fetchBibleBooks().then(function () {
      renderBible();
    });
  }

  function closeSheets() {
    dom.announcementOverlay.classList.add('hidden');
    dom.bibleOverlay.classList.add('hidden');
  }

  function gotoChapter(bookSn, chapter) {
    var book = bookBySn(bookSn);
    if (!book) {
      return;
    }
    if (chapter < 1) {
      // 退回上一卷最后一章；已是第一卷则停在原地。
      var books = state.bibleBooks || [];
      for (var i = 0; i < books.length; i++) {
        if (books[i].sn === bookSn && i > 0) {
          var prev = books[i - 1];
          state.bibleBookSn = prev.sn;
          state.bibleChapter = prev.chapterCount;
          state.bibleView = 'reading';
          renderBible();
          return;
        }
      }
      return;
    }
    if (chapter > book.chapterCount) {
      // 进入下一卷第一章；已是最后一卷则停在原地。
      var list = state.bibleBooks || [];
      for (var j = 0; j < list.length; j++) {
        if (list[j].sn === bookSn && j < list.length - 1) {
          var next = list[j + 1];
          state.bibleBookSn = next.sn;
          state.bibleChapter = 1;
          state.bibleView = 'reading';
          renderBible();
          return;
        }
      }
      return;
    }
    state.bibleChapter = chapter;
    state.bibleView = 'reading';
    renderBible();
  }

  function showToolRow() {
    state.presentable = true;
    dom.toolRow.classList.remove('hidden');
  }

  function hideToolRow() {
    state.presentable = false;
    dom.toolRow.classList.add('hidden');
    closeSheets();
  }

  /* ---- 观看统计：会话键 + 心跳 ---- */

  /* 会话键优先复用 sessionStorage 中的值：刷新页面不新增统计记录，
     新开标签页或关闭后重开才算一次新的观看。 */
  function ensureSessionKey() {
    try {
      var stored = window.sessionStorage.getItem(SESSION_STORAGE_KEY);
      if (stored) {
        return stored;
      }
      var generated = newSessionKey();
      window.sessionStorage.setItem(SESSION_STORAGE_KEY, generated);
      return generated;
    } catch (e) {
      /* 隐私模式下 sessionStorage 可能不可用：退化为每次页面驻留一个键。 */
      return state.memorySessionKey || (state.memorySessionKey = newSessionKey());
    }
  }

  function newSessionKey() {
    if (window.crypto && typeof window.crypto.randomUUID === 'function') {
      return window.crypto.randomUUID();
    }
    return 'sess-' + Date.now().toString(36) + '-' + Math.random().toString(36).slice(2, 10);
  }

  function startHeartbeat() {
    stopHeartbeat();
    sendHeartbeat();
    state.heartbeatTimer = window.setInterval(sendHeartbeat, HEARTBEAT_INTERVAL_MS);
  }

  function stopHeartbeat() {
    if (state.heartbeatTimer) {
      window.clearInterval(state.heartbeatTimer);
      state.heartbeatTimer = null;
    }
  }

  function sendHeartbeat() {
    var query = parseQuery();
    if (!query.room) {
      return;
    }
    var body = {
      roomCode: query.room,
      sessionKey: ensureSessionKey()
    };
    if (query.tenant !== '') {
      var tenantId = Number(query.tenant);
      if (!isNaN(tenantId)) {
        body.tenantId = tenantId;
      }
    }
    fetch(heartbeatEndpoint(), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
      cache: 'no-store',
      keepalive: true
    }).then(function (response) {
      return response.json();
    }).then(function (payload) {
      if (payload && payload.code === 0 && payload.data) {
        renderViewersBadge(Number(payload.data.onlineCount) || 0);
      }
    }).catch(function () {
      /* 心跳失败静默忽略：下一轮自动补发，不打扰观看体验。 */
    });
  }

  function renderViewersBadge(online) {
    if (online > 0) {
      dom.viewersCount.textContent = String(online);
      dom.viewersBadge.classList.remove('hidden');
    } else {
      dom.viewersBadge.classList.add('hidden');
    }
  }

  /* ---- 回放库 ---- */

  function fetchReplays() {
    return requestWithTimeout(replaysEndpoint())
      .then(function (response) {
        return response.json().then(function (payload) {
          return { status: response.status, payload: payload };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.status >= 200 && result.status < 300 && payload.code === 0 && payload.data) {
          return payload.data.list || [];
        }
        return [];
      })
      .catch(function () {
        return [];
      });
  }

  function renderReplays(items) {
    state.replayItems = items || [];
    if (!state.replayItems.length) {
      dom.replaySection.classList.add('hidden');
      dom.replayList.innerHTML = '';
      return;
    }
    var cards = state.replayItems.map(function (item, index) {
      return '<button type="button" role="listitem" class="replay-card" data-index="' + index + '">' +
        '<span class="replay-cover">' +
        '<svg viewBox="0 0 24 24" fill="none" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="8.6"/><path d="M10.2 9l5 3-5 3z"/></svg>' +
        '</span>' +
        '<span class="replay-info">' +
        '<span class="replay-title">' + escapeHtml(item.title || '往期回放') + '</span>' +
        '<span class="replay-date">' + escapeHtml(item.liveDate || '') + '</span>' +
        '</span>' +
        '</button>';
    }).join('');
    dom.replayList.innerHTML = cards;
    dom.replaySection.classList.remove('hidden');
  }

  /* 点击回放卡片切换播放源；当前正在播放的场次高亮并禁点。 */
  function bindReplayList() {
    dom.replayList.addEventListener('click', function (event) {
      var card = event.target.closest('.replay-card');
      if (!card || card.classList.contains('active')) {
        return;
      }
      var index = Number(card.getAttribute('data-index'));
      var item = state.replayItems[index];
      if (!item || !item.liveUrl) {
        return;
      }
      Array.prototype.forEach.call(dom.replayList.querySelectorAll('.replay-card'), function (node) {
        node.classList.remove('active');
      });
      card.classList.add('active');
      dom.playerError.classList.add('hidden');
      attachStream(item.liveUrl);
    });
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
    dom.viewersBadge.classList.add('hidden');
    dom.playerError.classList.add('hidden');
    dom.unmuteHint.classList.add('hidden');
    dom.liveRoom.classList.add('hidden');
    dom.coverDate.classList.add('hidden');
    dom.calendarBar.classList.add('hidden');
    dom.shareBar.classList.add('hidden');
    dom.replaySection.classList.add('hidden');
    dom.toolRow.classList.add('hidden');
    closeSheets();
    state.presentable = false;
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
    dom.calendarBar.classList.add('hidden');
    dom.viewersBadge.classList.add('hidden');
    dom.viewersCount.textContent = '';
    if (play.state === 2) {
      dom.replayBadge.classList.remove('hidden');
      setStateChip('replay', '已结束 · 回放');
    } else {
      dom.liveBadge.classList.remove('hidden');
      setStateChip('live', '直播中');
    }
    renderHeader(play);
    renderProgram(play);
    dom.shareBar.classList.remove('hidden');
    showToolRow();
    // 观看心跳只在可播状态（直播中 / 回放）发送；预告与无直播不产生观看记录。
    startHeartbeat();
    // 回放态补充展示往期回放列表；直播中不展示，避免干扰当前场次。
    if (play.state === 2) {
      fetchReplays().then(function (items) {
        renderReplays(items);
        highlightReplayCard(play.liveId);
      });
    } else {
      renderReplays([]);
    }
    attachStream(play.liveUrl);
  }

  /* 在回放列表中标记当前正在播放的场次。 */
  function highlightReplayCard(liveId) {
    if (!liveId) {
      return;
    }
    Array.prototype.forEach.call(dom.replayList.querySelectorAll('.replay-card'), function (node) {
      var index = Number(node.getAttribute('data-index'));
      var item = state.replayItems[index];
      if (item && item.liveId === liveId) {
        node.classList.add('active');
      }
    });
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
    // 预告态展示日历订阅入口；进行中 / 回放 / 无直播态保持隐藏。
    // 缺少房间参数时页面直接进入提示态，不会走到这里。
    dom.calendarBar.classList.remove('hidden');
    // 预告态没有可播内容：不发送观看心跳，仅保留分享入口。
    stopHeartbeat();
    dom.shareBar.classList.remove('hidden');
    showToolRow();
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
    stopHeartbeat();
    destroyPlayer();
    hideToolRow();
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

  /* 公告正文按换行渲染为段落；空行折叠，避免注入风险先转义再包裹。 */
  function renderMultiline(value) {
    var lines = String(value == null ? '' : value).split(/\r?\n/);
    return lines.map(function (line) {
      return '<span class="announcement-line">' + escapeHtml(line) + '</span>';
    }).join('');
  }

  /* ---- 日历订阅 ---- */

  /* 轻提示：居中底部浮现，2 秒后自动消失。 */
  function showToast(message) {
    if (state.toastTimer) {
      window.clearTimeout(state.toastTimer);
    }
    dom.toast.textContent = message;
    dom.toast.classList.add('visible');
    state.toastTimer = window.setTimeout(function () {
      dom.toast.classList.remove('visible');
      state.toastTimer = null;
    }, 2000);
  }

  /* 复制文本：优先异步剪贴板，非安全上下文回退到隐藏输入框方案。 */
  function copyText(text, onDone) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(function () {
        onDone(true);
      }, function () {
        onDone(legacyCopy(text));
      });
      return;
    }
    onDone(legacyCopy(text));
  }

  function legacyCopy(text) {
    var input = document.createElement('textarea');
    input.value = text;
    input.setAttribute('readonly', '');
    input.style.position = 'fixed';
    input.style.opacity = '0';
    document.body.appendChild(input);
    input.select();
    var copied = false;
    try {
      copied = document.execCommand('copy');
    } catch (e) {
      copied = false;
    }
    document.body.removeChild(input);
    return copied;
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
  bindReplayList();

  /* ---- 公告与圣经面板事件 ---- */

  dom.toolAnnouncement.addEventListener('click', openAnnouncementSheet);
  dom.toolBible.addEventListener('click', openBibleSheet);

  /* 遮罩与关闭按钮统一收回面板：data-sheet-close 标记关闭意图。 */
  Array.prototype.forEach.call(document.querySelectorAll('[data-sheet-close]'), function (node) {
    node.addEventListener('click', closeSheets);
  });

  /* 面板内容区统一事件委托：书卷 / 章节点击进入下一级，
     sheet-retry 按当前视图重试失败加载。 */
  dom.bibleBody.addEventListener('click', function (event) {
    var retry = event.target.closest('.sheet-retry');
    if (retry) {
      if (state.bibleView === 'reading' && state.bibleBookSn && state.bibleChapter) {
        /* 阅读视图原地重取当前章。 */
        dom.bibleBody.innerHTML = '<div class="sheet-empty">加载中…</div>';
        dom.bibleBody.scrollTop = 0;
        fetchBibleChapter(state.bibleBookSn, state.bibleChapter).then(function (data) {
          renderBibleVerses(data);
        });
        return;
      }
      state.bibleView = 'books';
      state.bibleBookSn = null;
      state.bibleChapter = null;
      state.bibleBooks = null;
      dom.bibleBody.innerHTML = '<div class="sheet-empty">加载中…</div>';
      fetchBibleBooks().then(function () {
        renderBible();
      });
      return;
    }
    var bookCard = event.target.closest('.bible-book');
    if (bookCard) {
      state.bibleBookSn = Number(bookCard.getAttribute('data-sn'));
      state.bibleView = 'chapters';
      renderBible();
      return;
    }
    var chapterCard = event.target.closest('.bible-chapter');
    if (chapterCard) {
      gotoChapter(state.bibleBookSn, Number(chapterCard.getAttribute('data-chapter')));
    }
  });

  /* 公告面板失败重试：与打开面板走同一加载路径。 */
  dom.announcementList.addEventListener('click', function (event) {
    if (!event.target.closest('.sheet-retry')) {
      return;
    }
    openAnnouncementSheet();
  });

  dom.bibleBack.addEventListener('click', function () {
    if (state.bibleView === 'reading') {
      state.bibleView = 'chapters';
    } else if (state.bibleView === 'chapters') {
      state.bibleView = 'books';
    }
    renderBible();
  });

  dom.biblePrev.addEventListener('click', function () {
    gotoChapter(state.bibleBookSn, (state.bibleChapter || 1) - 1);
  });

  dom.bibleNext.addEventListener('click', function () {
    gotoChapter(state.bibleBookSn, (state.bibleChapter || 1) + 1);
  });

  dom.calendarOpen.addEventListener('click', function () {
    /* 由用户环境决定能否唤起日历应用：能识别 ICS 的环境会直接打开。 */
    window.open(subscribeURL(), '_blank');
  });

  dom.calendarCopy.addEventListener('click', function () {
    copyText(subscribeURL(), function (copied) {
      showToast(copied ? '订阅链接已复制' : '复制失败，请长按地址栏复制');
    });
  });

  /* 分享入口：支持系统分享的环境唤起 navigator.share，
     其余环境降级为复制观播页链接，结果通过轻提示反馈。 */
  function handleShare() {
    var shareUrl = window.location.href;
    if (navigator.share) {
      navigator.share({
        title: document.title,
        url: shareUrl
      }).catch(function () { /* 用户取消分享不提示 */ });
      return;
    }
    copyText(shareUrl, function (copied) {
      showToast(copied ? '链接已复制' : '复制失败，请长按地址栏复制');
    });
  }

  dom.shareButton.addEventListener('click', handleShare);

  /* 桌面等无系统分享面板的环境：按钮文案降级为“复制链接”，
     与实际行为保持一致。 */
  if (!navigator.share) {
    var shareLabel = document.getElementById('share-button-label');
    if (shareLabel) {
      shareLabel.textContent = '复制链接';
    }
  }

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
     播放态不打断，交由 hls.js 自动回到直播边缘；播放态补发一次心跳，
     避免后台期间的超时间隔把观众从在线名单中挤掉。 */
  document.addEventListener('visibilitychange', function () {
    if (document.visibilityState !== 'visible') {
      return;
    }
    if (state.polling) {
      refresh();
    } else if (state.heartbeatTimer) {
      sendHeartbeat();
    }
  });

  /* 页面卸载前不再补发心跳：sendBeacon 不必要，60 秒窗口会自然过期。 */
  refresh();
})();
