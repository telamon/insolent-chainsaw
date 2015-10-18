var express = require('express'),
	fs = require('fs')
	app= express();

app.get('/',function(req,res){
	res.set('Content-Type', 'text/plain');
	var delay = req.query.delay || 0
	setTimeout(function(){
		res.status(200).send(fs.readFileSync('/version.txt'))
	},delay)
})
app.get('/exit',function(req,res){
	process.exit()
})
app.listen(7890)