<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
    <v-toolbar-title>Shopify Application</v-toolbar-title>

    </v-app-bar>

    <v-main>

 <v-card
    color="grey lighten-4"
  >
<v-container>
    <v-row>
      <v-col>
      <v-select
        v-model="selection"
        :items="items"
        label="Search By"
      ></v-select>
      </v-col>

  <v-col
    cols="12"
    md="8"
  >
      <v-text-field
        v-if="selection==='Text'"
        v-model="text"
        hide-details
        prepend-icon="mdi-magnify"
        single-line
      ></v-text-field>

    <v-combobox
      v-if="selection==='Tags'"
      v-model="tags"
      :items="tagOptions"
      hide-selected
      label="Add some tags"
      multiple
      persistent-hint
      small-chips
    >
      <template v-slot:no-data>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title>
              No results matching "
              <strong>{{ search }}</strong>". Press <kbd>enter</kbd> to create a new one
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-combobox>

      <v-file-input
        v-if="selection==='Image'"
        label="Image"
      ></v-file-input>
    </v-col>
    </v-row>
</v-container>

    </v-card>

<v-container>
        <v-row>
          <v-col
            v-for="n in 9"
            :key="n"
            class="d-flex child-flex"
            cols="2"
          >
            <v-img
              :src="`https://picsum.photos/500/300?image=${n * 5 + 10}`"
              :lazy-src="`https://picsum.photos/10/6?image=${n * 5 + 10}`"
              aspect-ratio="1"
              class="grey lighten-2"
            >
              <template v-slot:placeholder>
                <v-row
                  class="fill-height ma-0"
                  align="center"
                  justify="center"
                >
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
export default {
  name: 'App',

  data: () => ({
    selection: 'Tags',
    items: [
      'Image',
      'Text',
      'Tags',
    ],
    tagOptions: ['Gaming', 'Programming', 'Vue', 'Vuetify'],
    tags: [],
    search: null,
  }),
};
</script>
