<template>
  <div class="p-index">
    <client-only>
      <v-parallax dark :src="require('~/assets/image.jpg')">
        <v-container class="title">
          <v-row justify="center" align="center">
            <h1>Title</h1>
          </v-row>
        </v-container>
      </v-parallax>
    </client-only>
    <v-container>
      <v-row class="outline pa-4 pa-sm-8 pa-lg-16">
        居留地女の間で
        その晩、私は隣室のアレキサンダー君に案内されて、始めて横浜へ遊びに出かけた。
        アレキサンダー君は、そんな遊び場所に就いてなら、日本人の私なんぞよりも、遙かに詳かに心得ていた。
        アレキサンダー君は、その自ら名告るところに依れば、旧露国帝室付舞踏師で、革命後上海から日本へ渡って来たのだが、踊を以て生業とすることが出来なくなって、今では銀座裏の、西洋料理店某でセロを弾いていると云う、つまり街頭で、よく見かける羅紗売りより僅かばかり上等な類のコーカサス人である。
        それでも、遉にコーカサス生れの故か、髪も眼も真黒で却々眉目秀麗ハンサムな男だったので、貧乏なのにも拘らず、居留地女の間では、格別可愛がられているらしい。
        ――アレキサンダー君は、露西亜語の他に、拙い日本語と、同じ位拙い英語とを喋ることが出来る。
      </v-row>
      <v-row justify="center" align="center">
        <v-col cols="12">
          <client-only>
            <v-carousel>
              <v-carousel-item
                v-for="item in items"
                :key="item.id"
                :src="item.thumbnailUrl"
              />
            </v-carousel>
          </client-only>
        </v-col>
      </v-row>
      <v-row>
        <v-col
          v-for="item in items"
          :key="item.id"
          sm="6"
          lg="3"
          @click="showWorkModal(item)"
        >
          <v-hover v-slot="{ hover }">
            <v-card class="card">
              <v-img
                class="thumbnail"
                :class="{ 'on-hover': hover }"
                :src="item.thumbnailUrl"
              >
                <v-row class="fill-height mx-0" align="end">
                  <div>
                    <v-card-title>{{ item.title }}</v-card-title>
                    <v-card-subtitle>{{ item.author }}</v-card-subtitle>
                  </div>
                </v-row>
              </v-img>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>
    </v-container>
    <modal :is-visible.sync="localState.isWorkModalVisible">
      <v-btn
        v-if="currentIndex > 0"
        slot="previousButton"
        icon
        large
        dark
        @click="seek(-1)"
      >
        <v-icon> mdi-chevron-left </v-icon>
      </v-btn>
      <work slot="content" class="work" :content="localState.modalContent" />
      <v-btn
        v-if="currentIndex < countOfItems - 1"
        slot="nextButton"
        icon
        large
        dark
        @click="seek(1)"
      >
        <v-icon> mdi-chevron-right </v-icon>
      </v-btn>
    </modal>
  </div>
</template>

<style>
.thumbnail > .v-image__image {
  transition: opacity 0.3s ease-in-out;
}

.thumbnail:not(.on-hover) > .v-image__image {
  opacity: 0.7;
}
</style>

<style scoped>
.outline {
  text-align: center;
}

.thumbnail {
  cursor: pointer;
}

.work {
  height: 100%;
}
</style>

<script lang="ts">
import {
  computed,
  defineComponent,
  getCurrentInstance,
  reactive,
} from '@nuxtjs/composition-api'
import Work from '~/components/Work.vue'
import { Work as IWork } from '~/store/works'
import Modal from '~/components/Modal.vue'

export default defineComponent({
  components: {
    Work,
    Modal,
  },
  setup() {
    const self = getCurrentInstance()
    const items = computed(() => self?.$store.state.works.works)
    const currentIndex = computed(() =>
      self?.$store.state.works.works.indexOf(localState.modalContent)
    )
    const countOfItems = computed(() => self?.$store.state.works.works.length)
    const localState = reactive({
      isWorkModalVisible: false,
      modalContent: {
        id: '',
        author: '',
        title: '',
        description: '',
        thumbnailUrl: '',
        contentUrl: '',
        createdAt: new Date(),
        updatedAt: new Date(),
      } as IWork,
    })

    function showWorkModal(content: IWork) {
      localState.modalContent = content
      localState.isWorkModalVisible = true
    }

    function seek(delta: number) {
      localState.modalContent =
        self?.$store.state.works.works[currentIndex.value + delta]
    }

    return {
      items,
      currentIndex,
      countOfItems,
      localState,
      showWorkModal,
      seek,
    }
  },
})
</script>
