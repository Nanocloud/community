import Ember from 'ember';

export default Ember.Component.extend({

  actions: {
    toggleVdiWindow() {
      if (this.get('toggleWindow')) {
        this.toggleWindow();
      }
      else {
        this.toggleProperty('stateVisible');
      }
    },
  }
});
