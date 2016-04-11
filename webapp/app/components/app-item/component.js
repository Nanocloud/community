import Ember from 'ember';

export default Ember.Component.extend({

    application: null,
    isEditing: false,
    connectionName: null,
    showSingleTab: false,
    session: Ember.inject.service('session'),
    unpublishState: false,
    isUnpublished: false,
 
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
            this.toast.error("Application hasn't been renamed");
          });
      },

      cancelEditMode() {
        this.set('isEditing', false);
      },

      unpublish() {
        this.set('unpublishState', true);
        this.application.destroyRecord()
          .then(() => {
            this.toast.success("Application has been unpublished successfully");
            this.set('unpublishState', false);
            this.set('isUnpublished', true);
          }, () => {
            this.toast.error("Application hasn't been unpublished");
            this.set('unpublishState', false);
          });
      }
    }
});
