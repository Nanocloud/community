import Ember from 'ember';

export default Ember.Component.extend({

  that: this,
	token: null,
	loadingFile: null,

  showElement() {
    $('.element-active-state').css("opacity", "1");
  },

  hideElement() {
    $('.element-active-state').css("opacity", "0");
  },

  dragEnter(e) {
    this.showElement();
  },

  dragLeave(e) {
    this.hideElement();
  },

  drop() {
  },

  initDropZone: function(){

    var that = this.that;

		var flow = new Flow({
      target: 'http://localhost:8080/upload',
			headers: { Authorization: "Bearer YhcWIstQBdULmysDQXXkACayL" },
			singleFile: true
		});

		flow.assignDrop(this.element);

		flow.on('filesSubmitted', function(){
      console.log('submit file');
			flow.upload();
		});

		flow.on('complete', function(event, flow){
      console.log('complete');
		});

		flow.on('uploadStart', function(event, flow){
		});

		flow.on('fileProgress', function(flow, file){
      console.log(flow.progress());
		});

	}.on('didInsertElement')
});
