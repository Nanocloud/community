import Ember from 'ember';

export default Ember.Controller.extend({

  init: function() {
    this.model = this.store.createRecord('user', {});
  },

  actions: {
    add: function() {
      this.model.save().then(function() {
        this.set('errorMessage', "User successfully created");
      }.bind(this), function(errorMessage) {
        this.set('errorMessage', errorMessage);        
      }.bind(this));
    }
  }
});
