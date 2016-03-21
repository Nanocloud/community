import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  connectionName: null,

  actions: {
    toggleSingleTab(connectionName) {
      console.log('single tab with: ' + connectionName);
      this.set('connectionName', connectionName);
      this.toggleProperty('showSingleTab');
    }
  }
});
