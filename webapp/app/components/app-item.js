import Ember from 'ember';

export default Ember.Component.extend({

    application: null,
    isEditing: false,
    connectionName: null,
    showSingleTab: false,
    
    actions : {

      toggleSingleTab(connectionName) {
        console.log('single tab with: ' + connectionName);
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
          }, (errorMessage) => {
          });
      },

      cancelEditMode() {
        this.set('isEditing', false);
      }
    }
});
