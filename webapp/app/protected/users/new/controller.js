import Ember from 'ember';

export default Ember.Controller.extend({

  passwordConfirmation: null,

  actions: {
    add() {
      if (this.get('model.password') !== this.get('passwordConfirmation')) {
        this.toast.error('Password must match');
        return ;
      }
      this.model
        .validate()
        .then(({ m, validations }) => {

          if (validations.get('isInvalid') === true) {
            return this.toast.error('Cannot create user');
          }

          this.model.save()
            .then(() => {
              this.transitionToRoute('protected.users');
            }, (errorMessage) => {
              this.toast.error('Cannot create new user : ' + errorMessage);
            });
        });
    }
  }
});
