import Ember from 'ember';

export default Ember.Component.extend({

  classNames: ["card-component"],

  click() {
    this.sendAction();
  },

  actions: {
  },
});
