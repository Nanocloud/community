import Ember from 'ember';

export default Ember.Component.extend({
  isVisible: false,
  session: Ember.inject.service('session'),
  store: Ember.inject.service('store'),
  selectedFile: null,
  history: { data: [ "C:\\" ], offset:0 },

  files: Ember.computed(function() {
    return (this.get('model'));
  }),

  initialize: function() {

    this.loadDirectory();
   
  }.on('becameVisible'),

  selectFile(file) {
    this.set('selectedFile', file);
  },

  selectDir(dir) {
    this.get('history').offset = this.get('history').offset + 1;
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
    var offset = this.get('history').offset;
    this.get('history').data.splice(offset, this.get('history').data.length - offset);

    this.get('history').data.pushObject(folder);
    this.loadDirectory();
  },

  goBack: function() {
    if (this.get('history').offset <= 0) return;
    this.get('history').offset = this.get('history').offset - 1;
    this.loadDirectory();
  },

  goNext: function() {
    if (this.get('history').offset > this.get('history').data.length) return;
    this.get('history').offset = this.get('history').offset + 1;
    this.loadDirectory();
  },
    
  pathToString: function() {
    var history = this.get('history');
    var data = history.data;
    var path = "";
    for (var i = 0; i <= history.offset; i++) {
      path += data[i] + "\\";
    }
    return (path);
  },

  actions: {
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
