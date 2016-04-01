import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),

  actions: {
    toggleSingleTab(connectionName) {
      this.set('connectionName', connectionName);
      this.toggleProperty('showSingleTab');
    },

    toggleFileExplorer() {
      this.toggleProperty('showFileExplorer');
    },

    updateModel() {
      console.log('trigger update model');
      
      var that = this;
      //this.get('model').update();

      that.get('model').update();
      window.titi = this;
    }
  }
});
