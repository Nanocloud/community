import Ember from 'ember';

export default Ember.Component.extend({
  isVisible: false,
  session: Ember.inject.service('session'),
  store: Ember.inject.service('store'),
  selectedFile: null,
  history: [ "C:\\" ],
  history_offset: 0,

  files: Ember.computed(function() {
    return (this.get('model'));
  }),

  historyData: Ember.computed('history_offset', function() {
    console.log('recalculate');
    return (this.pathToArray());
  }),

  initialize: function() {

    this.loadDirectory();
   
  }.on('becameVisible'),

  selectFile(file) {
    this.set('selectedFile', file);
  },

  selectDir(dir) {
    this.incrementProperty('history_offset');
    this.goToDirectory(dir);
  },

  loadDirectory: function() {
    var path = this.pathToString();
    this.get('store').query('file', { filename: path }) 
      .then(function(response) {
        this.set('items', response);
      }.bind(this));
  },

  goToDirectory: function(folder) {

    // removing from current
    var offset = this.get('history_offset');
    this.get('history').splice(offset, this.get('history').length - offset);

    this.get('history').pushObject(folder);
    this.loadDirectory();
  },

  goBack: function() {
    if (this.get('history_offset') <= 0) return;
    this.decrementProperty('history_offset');
    this.loadDirectory();
  },

  goNext: function() {
    if ((this.get('history_offset')+1) >= this.get('history').length) return;
    this.incrementProperty('history_offset');
    this.loadDirectory();
  },

  pathToArray: function() {
    var data = this.get('history');
    var offset = this.get('history_offset');
    var path = [];
    for (var i = 0; i <= offset; i++) {
      path.pushObject(data[i]);
    }
    return (path);
  },
    
  pathToString: function() {
    var data = this.get('history');
    var offset = this.get('history_offset');
    var path = "";
    for (var i = 0; i <= offset; i++) {
      path += data[i] + "\\";
    }
    return (path);
  },

  actions: {

    moveOffset(offset) {
      console.log('moving offset');
      this.set('history_offset', offset);
      this.loadDirectory();
    },

    toggleFileExplorer() {
      this.toggleProperty('isVisible');
    },

    clickItem(item) {
      if (item.get('type') == 'directory') {
        this.selectDir(item.id);
        return;
      }

      this.selectFile(item);
    },

    clickNextBtn() {
      this.goNext();
    },

    clickPrevBtn() {
      this.goBack();
    },
  }
});
