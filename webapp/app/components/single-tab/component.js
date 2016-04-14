import Ember from 'ember';

export default Ember.Component.extend({

  isVisible: false,
  connectionName: null,
  vdiWindowVisible: false,

  actions: {
    toggleSingleTab() {
      this.sendAction('onClose', this.get('connectionName'));
      this.toggleProperty('isVisible');
    },

    toggleVdiWindow() {
      this.toggleProperty('vdiWindowVisible');
    }
  }
});
