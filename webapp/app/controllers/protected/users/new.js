import Ember from 'ember';

export default Ember.Controller.extend({

  passwordConfirmation: null,

  actions: {
    add() {
      if (this.model.get('password') != this.get('passwordConfirmation')) {
        this.set('errorMessage', "Password must match");
        return ;
      }
      this.model.save()
      .then(() => {
        this.set('errorMessage', "User successfully created");
        controller.set('model', this.store.createRecord('user', {}));
      }, (errorMessage) => {
        this.set('errorMessage', errorMessage);
      });
    }
  }
});
