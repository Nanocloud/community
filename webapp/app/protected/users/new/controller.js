import Ember from 'ember';

export default Ember.Controller.extend({

  passwordConfirmation: null,
  loadState: false,

  actions: {
    add() {
      this.model
        .validate()
        .then(({ m, validations }) => {

          if (validations.get('isInvalid') === true) {
            return this.toast.error('Cannot create user');
          }

          this.set('loadState', true);
          this.model.save()
            .then(() => {
              this.set('loadState', false);
              this.transitionToRoute('protected.users');
            }, (errorMessage) => {
              this.set('loadState', false);
              this.toast.error('Cannot create new user : ' + errorMessage);
            });
        });
    }
  }
});
