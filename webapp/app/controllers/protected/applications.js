import Ember from 'ember';

export default Ember.Controller.extend({

  showSingleTab: false,
  actions: {
    toggleSingleTab() {
      this.set('connectionName', 'hapticPowershell');
      this.toggleProperty('showSingleTab');
    }
  }
});
