<template>
  <v-app id="inspire">
    <h1>Upload a File</h1>
    <form enctype="multipart/form-data">
      <input
        type="file"
        name="file"
        v-on:change="fileChange($event.target.files)"
      />
      <v-btn v-on:click="upload()">Upload</v-btn>
    </form>
    <v-card>
      <v-card-title>
        Traces
        <v-spacer></v-spacer>
        <v-text-field
          v-model="search"
          append-icon="search"
          label="Search"
          single-line
          hide-details
        ></v-text-field>
      </v-card-title>
      <v-data-table :headers="headers" :items="hostData" :search="search">
        <template v-slot:items="props">
          <tr
            @click.stop="
              getPortData(props.item.hostId, props.item.hostname);
              dialog = true;
            "
          >
            <td>{{ props.item.hostId }}</td>
            <td>{{ props.item.hostname }}</td>
            <td>{{ props.item.addr }}</td>
            <td>{{ props.item.addrType }}</td>
            <td>{{ props.item.updatedAt }}</td>
          </tr>
        </template>
        <v-alert v-slot:no-results :value="true" color="error" icon="warning"
          >Your search for "{{ search }}" found no results.</v-alert
        >
      </v-data-table>
    </v-card>

    <v-footer color="default" app>
      <span class="white--text">&copy; 2017</span>
    </v-footer>

    <v-dialog v-model="dialog">
      <v-card>
        <v-card-title class="headline">Details</v-card-title>
        <v-card-text>Data for Host {{ hostname }}</v-card-text>
        <v-data-table :headers="headersPorts" :items="portData">
          <template v-slot:items="props">
            <tr>
              <td>{{ props.item.protocol }}</td>
              <td>{{ props.item.portId }}</td>
              <td>{{ props.item.state }}</td>
              <td>{{ props.item.reason }}</td>
              <td>{{ props.item.name }}</td>
            </tr>
          </template>
        </v-data-table>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" flat="flat" @click="dialog = false"
            >Close</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import axios from "axios";
export default {
  data: () => ({
    files: new FormData(),
    search: "",
    dialog: false,
    hostname: "",
    headers: [
      { text: "Host ID", value: "hostId" },
      { text: "Hostname", value: "hostname" },
      { text: "IP Address", value: "addr" },
      { text: "Address Type", value: "addrType" },
      { text: "Updated At", value: "updatedAt" }
    ],
    headersPorts: [
      { text: "Host ID", value: "hostId" },
      { text: "Hostname", value: "hostname" },
      { text: "IP Address", value: "addr" },
      { text: "Address Type", value: "addrType" },
      { text: "Updated At", value: "updatedAt" }
    ],
    hostData: [],
    portData: []
  }),
  props: {
    source: String
  },
  methods: {
    async getPortData(hostId, hostname) {
      let res = await axios.get(
        `http://localhost:3000/v1/getports/${hostId}`,
        this.files,
        {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        }
      );
      this.hostname = hostname;
      this.portData = res.data;
    },
    fileChange(fileList) {
      this.files.append("file", fileList[0], fileList[0].name);
    },
    async upload() {
      // TODO: error handling
      await axios.post("http://localhost:3000/v1/upload.xml", this.files, {
        headers: {
          "Content-Type": "multipart/form-data"
        }
      });
      await this.refresh();
      this.files = new FormData();
    },
    async refresh() {
      let res = await axios.get(
        `http://localhost:3000/v1/gethosts`,
        this.files,
        {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        }
      );
      if (res.data !== null) {
        this.hostData = res.data;
      }
    }
  },
  async created() {
    await this.refresh();
  }
};
</script>
