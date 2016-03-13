import Ember from 'ember';

export default Ember.Controller.extend({

  passwordConfirmation: null,

  actions: {
    add() {
      if (this.get('model.password') !== this.get('passwordConfirmation')) {
        this.set('errorMessage', "Password must match");
        return ;
      }
      this.model.save()
        .then(() => {
          this.transitionToRoute('protected.users');
      }, (errorMessage) => {
        this.set('errorMessage', errorMessage);
      });
    }
  }
});
