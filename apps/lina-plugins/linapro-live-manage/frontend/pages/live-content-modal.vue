<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  Button,
  DatePicker,
  Form,
  FormItem,
  Input,
  message,
  RadioGroup,
  Select,
} from 'ant-design-vue';

import { $t } from '#/locales';
import { TiptapEditor } from '#/components/tiptap';
import {
  liveAdd,
  liveInfo,
  liveUpdate,
  liveroomOptions,
  type LiveSong,
} from './live-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface SongRow {
  name: string;
  singer: string;
}

interface FormData {
  id?: number;
  roomId?: number;
  title: string;
  content: string;
  liveDate: string;
  pageUrl: string;
  pushUrl: string;
  liveUrl: string;
  coverUrl: string;
  songName: string;
  songRows: SongRow[];
  leadSinger: string;
  accompaniment: string;
  host: string;
  sermonTitle: string;
  preacher: string;
  preacherIdentity: string;
  scriptureRef: string;
  scriptureContent: string;
  outline: string;
  deviceInfo: string;
  reception: string;
  state: number;
  isPublic: number;
  startTime: number | null;
}

const defaultValues: FormData = {
  id: undefined,
  roomId: undefined,
  title: '',
  content: '',
  liveDate: '',
  pageUrl: '',
  pushUrl: '',
  liveUrl: '',
  coverUrl: '',
  songName: '',
  songRows: [],
  leadSinger: '',
  accompaniment: '',
  host: '',
  sermonTitle: '',
  preacher: '',
  preacherIdentity: '',
  scriptureRef: '',
  scriptureContent: '',
  outline: '',
  deviceInfo: '',
  reception: '',
  state: 0,
  isPublic: 1,
  startTime: null,
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-live-manage.drawer.liveEditTitle')
    : $t('plugin.linapro-live-manage.drawer.liveCreateTitle'),
);

const formRules = reactive({
  roomId: [
    { message: $t('plugin.linapro-live-manage.validation.roomId'), required: true },
  ],
  title: [{ message: $t('plugin.linapro-live-manage.validation.title'), required: true }],
  state: [{ message: $t('plugin.linapro-live-manage.validation.state'), required: true }],
  isPublic: [
    { message: $t('plugin.linapro-live-manage.validation.isPublic'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const dictStore = useDictStore();
const roomOptions = ref<{ id: number; roomName: string; roomCode: string }[]>([]);
const liveStateOptions = ref<{ label: string; value: number }[]>([]);
const livePublicOptions = ref<{ label: string; value: number }[]>([]);

async function loadRoomOptions(keyword?: string) {
  const res = await liveroomOptions(keyword);
  roomOptions.value = res?.list ?? [];
}

function addSongRow() {
  formData.value.songRows.push({ name: '', singer: '' });
}

function removeSongRow(index: number) {
  formData.value.songRows.splice(index, 1);
}

function serializeSongRows(): string {
  const songs: LiveSong[] = formData.value.songRows
    .map((row, index) => ({
      name: row.name.trim(),
      singer: row.singer.trim(),
      order: index + 1,
    }))
    .filter((song) => song.name !== '' || song.singer !== '');
  if (songs.length === 0) {
    return '';
  }
  return JSON.stringify(songs);
}

function parseSongRows(songList: string): SongRow[] {
  if (!songList) {
    return [];
  }
  try {
    const parsed = JSON.parse(songList) as Array<{
      name?: string;
      singer?: string;
    }>;
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed.map((song) => ({
      name: song?.name ?? '',
      singer: song?.singer ?? '',
    }));
  } catch {
    return [];
  }
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[960px]',
  fullscreenButton: true,
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const [states, publics] = await Promise.all([
      dictStore.getDictOptionsAsync('plugin_live_state'),
      dictStore.getDictOptionsAsync('plugin_live_public'),
    ]);
    liveStateOptions.value = states.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    livePublicOptions.value = publics.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    await loadRoomOptions();
    const data = modalApi.getData();
    if (data?.id) {
      modalApi.setState({ confirmLoading: true });
      try {
        const record = await liveInfo(data.id);
        formData.value = {
          id: record.id,
          roomId: record.roomId,
          title: record.title,
          content: record.content || '',
          liveDate: record.liveDate || '',
          pageUrl: record.pageUrl || '',
          pushUrl: record.pushUrl || '',
          liveUrl: record.liveUrl || '',
          coverUrl: record.coverUrl || '',
          songName: record.songName || '',
          songRows: parseSongRows(record.songList),
          leadSinger: record.leadSinger || '',
          accompaniment: record.accompaniment || '',
          host: record.host || '',
          sermonTitle: record.sermonTitle || '',
          preacher: record.preacher || '',
          preacherIdentity: record.preacherIdentity || '',
          scriptureRef: record.scriptureRef || '',
          scriptureContent: record.scriptureContent || '',
          outline: record.outline || '',
          deviceInfo: record.deviceInfo || '',
          reception: record.reception || '',
          state: record.state,
          isPublic: record.isPublic,
          startTime: record.startTime ?? null,
        };
      } finally {
        modalApi.setState({ confirmLoading: false });
      }
    } else {
      formData.value = { ...defaultValues };
      resetFields();
    }
  },
});

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();

    const { id, songRows: _songRows, ...values } = formData.value;
    const submitData = {
      ...values,
      songList: serializeSongRows(),
    };
    if (isEdit.value && id) {
      await liveUpdate(id, submitData);
      message.success($t('pages.common.updateSuccess'));
    } else {
      await liveAdd(submitData);
      message.success($t('pages.common.createSuccess'));
    }
    emit('reload');
    modalApi.close();
  } catch (error) {
    console.error(error);
  } finally {
    modalApi.lock(false);
  }
}
</script>

<template>
  <Modal :title="title">
    <Form layout="vertical">
      <FormItem
        :label="$t('plugin.linapro-live-manage.fields.roomId')"
        v-bind="validateInfos.roomId"
      >
        <Select
          v-model:value="formData.roomId"
          :filterOption="false"
          :options="
            roomOptions.map((room) => ({
              label: `${room.roomName}（${room.roomCode}）`,
              value: room.id,
            }))
          "
          :placeholder="$t('plugin.linapro-live-manage.placeholders.roomId')"
          :disabled="!!formData.id"
          show-search
          @search="loadRoomOptions"
        />
      </FormItem>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.title')"
          v-bind="validateInfos.title"
        >
          <Input
            v-model:value="formData.title"
            :maxlength="512"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.title')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.liveDate')">
          <DatePicker
            v-model:value="formData.liveDate"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.liveDate')"
            class="w-full"
            value-format="YYYY-MM-DD"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.state')"
          v-bind="validateInfos.state"
        >
          <RadioGroup
            v-model:value="formData.state"
            button-style="solid"
            option-type="button"
            :options="liveStateOptions"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.isPublic')"
          v-bind="validateInfos.isPublic"
        >
          <RadioGroup
            v-model:value="formData.isPublic"
            button-style="solid"
            option-type="button"
            :options="livePublicOptions"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.startTime')">
          <DatePicker
            v-model:value="formData.startTime"
            class="w-full"
            show-time
            :placeholder="$t('plugin.linapro-live-manage.placeholders.startTime')"
            value-format="x"
          />
        </FormItem>
      </div>

      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-live-manage.fields.liveUrl')">
          <Input
            v-model:value="formData.liveUrl"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.liveUrl')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.pushUrl')">
          <Input
            v-model:value="formData.pushUrl"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.pushUrl')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.pageUrl')">
          <Input
            v-model:value="formData.pageUrl"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.pageUrl')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.coverUrl')">
          <Input
            v-model:value="formData.coverUrl"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.coverUrl')"
          />
        </FormItem>
      </div>

      <FormItem :label="$t('plugin.linapro-live-manage.fields.content')">
        <TiptapEditor v-model="formData.content" :height="240" scene="live_content" />
      </FormItem>

      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-live-manage.fields.songName')">
          <Input v-model:value="formData.songName" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.leadSinger')">
          <Input v-model:value="formData.leadSinger" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.accompaniment')">
          <Input v-model:value="formData.accompaniment" />
        </FormItem>
      </div>

      <FormItem :label="$t('plugin.linapro-live-manage.fields.songList')">
        <div class="flex flex-col gap-2">
          <div
            v-for="(row, index) in formData.songRows"
            :key="index"
            class="flex items-center gap-2"
          >
            <span class="w-8 text-center">{{ index + 1 }}</span>
            <Input
              v-model:value="row.name"
              :maxlength="512"
              :placeholder="$t('plugin.linapro-live-manage.placeholders.songName')"
              class="flex-1"
            />
            <Input
              v-model:value="row.singer"
              :maxlength="256"
              :placeholder="$t('plugin.linapro-live-manage.placeholders.singer')"
              class="flex-1"
            />
            <Button danger size="small" @click="removeSongRow(index)">
              {{ $t('pages.common.delete') }}
            </Button>
          </div>
          <Button class="w-fit" size="small" type="dashed" @click="addSongRow">
            {{ $t('plugin.linapro-live-manage.messages.addSong') }}
          </Button>
        </div>
      </FormItem>

      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-live-manage.fields.host')">
          <Input v-model:value="formData.host" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.sermonTitle')">
          <Input v-model:value="formData.sermonTitle" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.preacher')">
          <Input v-model:value="formData.preacher" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.preacherIdentity')">
          <Input v-model:value="formData.preacherIdentity" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.scriptureRef')">
          <Input v-model:value="formData.scriptureRef" />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-live-manage.fields.reception')">
          <Input v-model:value="formData.reception" />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-live-manage.fields.scriptureContent')">
        <Input.TextArea v-model:value="formData.scriptureContent" :rows="3" />
      </FormItem>
      <FormItem :label="$t('plugin.linapro-live-manage.fields.outline')">
        <Input.TextArea v-model:value="formData.outline" :rows="3" />
      </FormItem>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-live-manage.fields.deviceInfo')">
          <Input v-model:value="formData.deviceInfo" :maxlength="512" />
        </FormItem>
      </div>
    </Form>
  </Modal>
</template>
