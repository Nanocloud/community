import Ember from 'ember';

export default Ember.Controller.extend({

  editingPassword: false,
  passwordConfirmation: null,
  errorMessage: null,

  actions: {

    removeDone: function() {
      this.get('model').destroyRecord();
      this.transitionToRoute('protected.users');
    },

    toggleEditPassword: function() {
      this.toggleProperty('editingPassword');
      this.set('model.password', "");
      this.set('passwordConfirmation', "");
      this.set('errorMessage', null);
    },

    changePassword: function() {
      if (this.get('model.password') !== this.get('passwordConfirmation')) {
        this.set('errorMessage', "Password must match");
        return ;
      }
      this.model.save()
        .then(() => {
          this.send('toggleEditPassword');
        }, (errorMessage) => {
          this.set('errorMessage', errorMessage);
        });
    }
  }
});
