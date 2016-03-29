import Ember from 'ember';

export default Ember.Controller.extend({

  showSingleTab: false,
  actions: {
    toggleSingleTab(name) {
      this.set('connectionName', name);
      this.toggleProperty('showSingleTab');
    }
  }
});
