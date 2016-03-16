import Ember from 'ember';

export default Ember.Component.extend({

  isVisible: false,

  actions: {
    toggleSingleTab() {
      this.toggleProperty('isVisible');
    }
  }
});
