import Ember from 'ember';

export default Ember.Component.extend({

  classNames: ["card-component"],

  click() {
    console.log("FICLICK"); 
    this.sendAction();
  },

  actions: {
  },
});
