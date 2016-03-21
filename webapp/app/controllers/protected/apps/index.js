import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  connectionName: null,

  actions: {
    toggleSingleTab(connectionName) {
      this.toggleProperty('showSingleTab');
    }
  }
});
