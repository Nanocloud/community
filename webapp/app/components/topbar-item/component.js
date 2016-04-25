import Ember from 'ember';

export default Ember.Component.extend({

  classNameBindings: ['stateEnabled'],

  actions: {
    clickAction() {
      this.sendAction("click");
    }
  }
});
