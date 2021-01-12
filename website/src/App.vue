<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <v-toolbar-title>Shopify Application</v-toolbar-title>
    </v-app-bar>

    <v-main>
      <v-card color="grey lighten-4">
        <v-container>
          <v-row>
            <v-col>
              <v-select
                v-model="selection"
                :items="items"
                label="Search By"
              ></v-select>
            </v-col>

            <v-col cols="12" md="8">
              <v-text-field
                v-if="selection === 'Text'"
                v-model="searchText"
                hide-details
                prepend-icon="mdi-magnify"
                single-line
                @keyup="searchByText"
              ></v-text-field>

              <v-combobox
                v-if="selection === 'Tags'"
                v-model="searchTags"
                :items="tagOptions"
                hide-selected
                label="Add some tags"
                multiple
                persistent-hint
                small-chips
                @change="searchByTags"
              >
                <template v-slot:no-data>
                  <v-list-item>
                    <v-list-item-content>
                      <v-list-item-title>
                        No results matching "
                        <strong>{{ search }}</strong
                        >". Press <kbd>enter</kbd> to create a new one
                      </v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>
                </template>
              </v-combobox>

              <v-file-input
                v-model="searchImage"
                v-if="selection === 'Image'"
                label="Image"
                @change="searchByImage"
              ></v-file-input>
            </v-col>
          </v-row>
        </v-container>
      </v-card>

      <v-container>
        <v-row>
          <v-col
            v-for="(imageUrl, i) in imageUrls"
            :key="i"
            class="d-flex child-flex"
            cols="2"
          >
            <v-img :src="imageUrl" aspect-ratio="1" class="grey lighten-2">
              <template v-slot:placeholder>
                <v-row class="fill-height ma-0" align="center" justify="center">
                  <v-progress-circular
                    indeterminate
                    color="grey lighten-5"
                  ></v-progress-circular>
                </v-row>
              </template>
            </v-img>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',

  data: () => ({
    selection: 'Tags',
    items: ['Image', 'Text', 'Tags'],
    tagOptions: ['Gaming', 'Programming', 'Vue', 'Vuetify'],
    search: null,

    searchText: '',
    searchTags: [],
    searchImage: null,
  }),
  computed: {
    imageUrls() {
      return [...Array(9).keys()].map(
        (n) => `https://picsum.photos/500/300?image=${n * 5 + 10}`,
      );
    },
    searchClient() {
      return axios.create({
        baseURL: process.env.VUE_APP_SEARCH_API_BASE_URL,
      });
    },
  },
  methods: {
    async searchByText() {
      // this.imageUrls = await this.searchClient.post('_search')
      console.log(this.searchText);
    },
    async searchByTags() {
      // this.imageUrls = await this.searchClient.post('_search')
      console.log(this.searchTags);
    },
    async searchByImage() {
      // this.imageUrls = await this.searchClient.post('_search')
      console.log(this.searchImage);
    },
  },
};
</script>
