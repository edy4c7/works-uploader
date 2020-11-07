<template>
  <div class="p-index">
    <v-parallax dark :src="require('~/assets/image.jpg')" class="fill-height">
      <v-container>
        <v-row justify="center" align="center">
          <h1>Title</h1>
        </v-row>
      </v-container>
    </v-parallax>
    <v-container>
      <v-row class="caption pa-4 pa-sm-8 pa-md-16">
        居留地女の間で
        その晩、私は隣室のアレキサンダー君に案内されて、始めて横浜へ遊びに出かけた。
        アレキサンダー君は、そんな遊び場所に就いてなら、日本人の私なんぞよりも、遙かに詳かに心得ていた。
        アレキサンダー君は、その自ら名告るところに依れば、旧露国帝室付舞踏師で、革命後上海から日本へ渡って来たのだが、踊を以て生業とすることが出来なくなって、今では銀座裏の、西洋料理店某でセロを弾いていると云う、つまり街頭で、よく見かける羅紗売りより僅かばかり上等な類のコーカサス人である。
        それでも、遉にコーカサス生れの故か、髪も眼も真黒で却々眉目秀麗ハンサムな男だったので、貧乏なのにも拘らず、居留地女の間では、格別可愛がられているらしい。
        ――アレキサンダー君は、露西亜語の他に、拙い日本語と、同じ位拙い英語とを喋ることが出来る。
      </v-row>
      <v-row justify="center" align="center">
        <v-carousel>
          <v-carousel-item v-for="item in items" :key="item">
            <v-row class="fill-height" align="center" justify="center">
              <v-img :src="require('~/assets/image.jpg')" />
            </v-row>
          </v-carousel-item>
        </v-carousel>
      </v-row>
      <v-row>
        <v-col v-for="item in items" :key="item" sm="6" md="3">
          <v-img :src="require('~/assets/image.jpg')" />
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<style>
.caption {
  text-align: center;
}
</style>

<script lang="ts">
import { defineComponent, ref } from '@vue/composition-api'
import Logo from '~/components/Logo.vue'
import VuetifyLogo from '~/components/VuetifyLogo.vue'

export default defineComponent({
  components: {
    Logo,
    VuetifyLogo,
  },

  setup() {
    const items = ref<string[]>(['item1', 'item2', 'item3', 'item4', 'item5'])

    return {
      items,
    }
  },
})
</script>
