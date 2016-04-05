import Ember from 'ember';

export default Ember.Component.extend({

    application: null,
    isEditing: false,
    connectionName: null,
    showSingleTab: false,
    
    actions : {

      toggleSingleTab(connectionName) {
        this.set('connectionName', connectionName);
        this.toggleProperty('showSingleTab');
      },

      toggleEditName() {
        this.toggleProperty('isEditing');
      },

      submitEditName() {
        this.toggleProperty('isEditing');
        this.application.save()
          .then(() => {
            this.toast.success("Application has been renamed successfully");
          }, () => {
            this.toast.success("Application hasn't been renamed");
          });
      },

      cancelEditMode() {
        this.set('isEditing', false);
      },

      unpublish() {
        this.application.destroyRecord();
      }
    }
});
