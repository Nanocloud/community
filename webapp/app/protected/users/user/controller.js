import Ember from 'ember';

export default Ember.Controller.extend({

  passwordConfirmation: null,

  actions: {

    removeDone: function() {
      this.get('model').destroyRecord();
      this.transitionToRoute('protected.users');
    },

    toggleEditPassword: function() {
      this.set('model.password', "");
      this.set('passwordConfirmation', "");
    },

    changePassword: function(defer) {

      this.get('model')
        .validate({ on: ['password'] })
        .then(({ m, validations }) => {

          if (validations.get('isInvalid') == true) {
            this.toast.error(this.get('model.validations.attrs.password.messages'));
            return defer.reject(this.get('model.validations.attrs.password.messages'));
          }

          if (this.get('model.password') !== this.get('passwordConfirmation')) {
            this.toast.error("Password doesn't match confirmation");
            return defer.reject("Password doesn't match confirmation");
          }

          this.model.save()
            .then(() => {
              defer.resolve();
              this.send('toggleEditPassword');
              this.toast.success('Password has been updated successfully');
            }, () => {
              defer.reject();
              this.toast.error("Password hasn't been updated");
            });
        });
    }
  }
});
