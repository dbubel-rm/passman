<template>
  <v-app>
    <v-toolbar app>
      <v-toolbar-title class="headline text-uppercase">
        <span class="font-weight-light">BF-Trace Viewer</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-content>
      <Traces />
    </v-content>
  </v-app>
</template>

<script>
import Traces from "./components/Traces";

export default {
  name: "App",
  components: {
    Traces
  },
  data() {
    return {
      //
    };
  }
};
</script>
