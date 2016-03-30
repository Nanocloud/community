import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,

  actions: {
    toggleSingleTab(connectionName) {
      this.set('connectionName', connectionName);
      this.toggleProperty('showSingleTab');
    },

    toggleFileExplorer() {
      this.toggleProperty('showFileExplorer');
    },
  }
});
