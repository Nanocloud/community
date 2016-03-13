import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    add() {
      this.model.save()
      .then(() => {
        this.set('errorMessage', "User successfully created");
      }, (errorMessage) => {
        this.set('errorMessage', errorMessage);
      });
    }
  }
});
