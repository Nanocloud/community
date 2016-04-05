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

    changePassword: function() {
      if (this.get('model.password') !== this.get('passwordConfirmation')) {
        this.toast.error("Password doesn't match confirmation");
        return ;
      }
      this.model.save()
        .then(() => {
          this.send('toggleEditPassword');
          this.toast.success('Password has been updated successfully');
        }, () => {
          this.toast.error("Password hasn't been updated");
        });
    }
  }
});
