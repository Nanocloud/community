import Ember from 'ember';

export default Ember.Component.extend({

  actions: {
    toggleVdiWindow() {
      this.toggleProperty('stateVisible');
    },

    clearList() {
      this.sendAction('flushHistory');
    }
  }
});
