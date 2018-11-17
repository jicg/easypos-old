(function($){
	Date.prototype.format = function(format) {  
	    /* 
	     * 使用例子:format="yyyy-MM-dd hh:mm:ss"; 
	     */  
	    var o = {  
	        "M+" : this.getMonth() + 1, // month  
	        "d+" : this.getDate(), // day  
	        "h+" : this.getHours(), // hour  
	        "m+" : this.getMinutes(), // minute  
	        "s+" : this.getSeconds(), // second  
	        "q+" : Math.floor((this.getMonth() + 3) / 3), // quarter  
	        "S" : this.getMilliseconds()  
	        // millisecond  
	    }  
	    
	    if (/(y+)/.test(format)) {  
	        format = format.replace(RegExp.$1, (this.getFullYear() + "").substr(4  
	                        - RegExp.$1.length));  
	    }  
	    
	    for (var k in o) {  
	        if (new RegExp("(" + k + ")").test(format)) {  
	            format = format.replace(RegExp.$1, RegExp.$1.length == 1  
	                            ? o[k]  
	                            : ("00" + o[k]).substr(("" + o[k]).length));  
	        }  
	    }  
	    return format;  
	};
	var Pos = function(){
		$.ajaxSetup({cache:false});
		this.templatehtml =  $("#item-body-template").html();
		this.el_orderno = $("#orderno");
		this.el_productno = $("#productno");
		this.el_customno = $("#customno");
		this.el_totamt = $("#totamt");
		this.el_trueamt = $("#trueamt");
		this.el_payamt = $("#payamt");
		this.el_desc = $("#desc");
		this.el_span_trueamt = $("#span-trueamt");
		this.el_span_payamt = $("#span-payamt");
		this.el_span_retamt = $("#span-retamt");

		this.el_item_body = $(".item-body");

		this.posdata = {
			orderno :"",
			customno :"",
			totamt:0,
			trueamt:0,
			payamt:0,
			retamt:0,
			desc:"",
			items:new Array()
		};
		this.tot_changed = false;
		this.item_changed = false;
	};

	Pos.prototype = {
		init:function(){
			this.loadOrderno();
			this.bind();
		},


		bind:function(){
			that = this ;
			this.el_productno.keydown(function(e) {
				if(e.keyCode==13){that.loadPdt($(this).val());}
			});
			this.el_customno.blur(function(e){
				that.posdata.customno = $(this).val();
			});
			this.el_desc.blur(function(e){
				that.posdata.desc = $(this).val();
			});
			this.el_payamt.blur(function(e){
				that.posdata.payamt = parseFloat($(this).val())||0;
				that.tot_changed = true;
				that.drawHtml();
			})
			this.el_trueamt.blur(function(e){
				that.posdata.trueamt = parseFloat($(this).val())||0;
				that.tot_changed = true;
				that.drawHtml();
			})
            this.el_productno.autocomplete({
                serviceUrl: '/pos/qpros/'+this.el_productno.val(),
                orientation:"top",
                formatResult:function(suggestion, currentValue) {
                	var val = suggestion.name+"	"+suggestion.value;
                    if (!currentValue) {
                        return val;
                    }
                    var pattern = '(' + $.Autocomplete.utils.escapeRegExChars(currentValue) + ')';
                    return val.replace(new RegExp(pattern,'gi'), '<strong>$1<\/strong>').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/&lt;(\/?strong)&gt;/g, '<$1>');
                },
                deferRequestBy:200
            });
		},

		loadOrderno:function(){
			that = this ;
			$.get('/pos/getno',function(data){
				that.el_orderno.val(data.data);
				that.posdata.orderno = data.data;
			},"json");
		},

		loadPdt:function(no){
			that = this;
			if(!no){
				that.error('请输入商品编号', '请重新输入！');
				return;
			}else{
				$.get('/pos/pro/'+no,function(data){
					if(data.code==0){
                        that.el_productno.val("");
						that.insertPdt(data.data);
					}else{
						that.error(data.msg);
					}
				},"json");
			}
		},

		insertPdt:function(obj){
			if(!obj.desc){
				obj.desc = obj.no ;
			}
			var ishas = false;
			for(var i =0;i<this.posdata.items.length ; i++){
				var tobj = this.posdata.items[i];
				if(tobj.product_no==obj.no){
					tobj.qty=tobj.qty+1;
					ishas = true;
					break;
				}
			}
			if(!ishas){
				var nobj = {};
				nobj.qty =1;
				nobj.product_id = obj.id;
				nobj.product_no = obj.no;
				nobj.product_desc = obj.desc;
				nobj.saleprice = obj.saleprice;
				nobj.trueprice = obj.saleprice;
				this.posdata.items.push(nobj);
			}
			this.item_changed = true;
			this.drawHtml();
		},

		drawHtml:function(){
			if(this.item_changed){
				this.updateItemData();
				var html = "";
				for(var i =0;i<this.posdata.items.length;i++){
					html = html+this.getSingleItem(this.posdata.items[i]);
				}
				this.el_item_body.html(html);
				this.tot_changed = true;
			}

			if(this.tot_changed){
				this.updateTotData();
				this.el_customno.val(this.posdata.customno);
				this.el_desc.val(this.posdata.desc);
				this.el_totamt.val(this.posdata.totamt);
				this.el_trueamt.val(this.posdata.trueamt);
				this.el_payamt.val(this.posdata.payamt);
				this.el_span_trueamt.html(this.posdata.trueamt);
				this.el_span_payamt.html(this.posdata.payamt);
				this.el_span_retamt.html(this.posdata.retamt);
			}

			this.item_changed =false;
			this.tot_changed =false;
		},

		getSingleItem:function(obj){
			return nano(this.templatehtml,obj)
		},

		updateItemData:function(){
			var totamt=0;
			for(var i =0;i<this.posdata.items.length ; i++){
				var tobj = this.posdata.items[i];
				tobj.amt = parseFloat((tobj.qty*tobj.trueprice).toFixed(2))
				totamt = totamt+tobj.amt
			}
			this.posdata.totamt = totamt;
			this.posdata.trueamt = totamt;
		},

		updateTotData:function(){
			var data = this.posdata;
			data.retamt = data.payamt-data.trueamt ;
			if(data.retamt<0){data.retamt=0;}
		},

		updatePdt:function(no,feild,e){
			var items = this.posdata.items ;
			for(var i =0;i<items.length;i++){
				if(items[i].product_no==no){
					items[i][feild] = (parseFloat($(e).val())||0.0);
					break;
				}
			}
			this.item_changed = true;
			this.drawHtml();
		},

		removePdt:function(no){
			var items = this.posdata.items ;
			for(var i =0;i<items.length;i++){
				if(items[i].product_no==no){
					items.splice(i,1); 
					break;
				}
			}
			this.item_changed = true;
			this.drawHtml();
		},

		createOrder:function(e){
			var that = this;
			var $btn = $(e).button('loading');
			var param = {data:JSON.stringify(this.posdata)};
			$.post("/pos/create",param,function(data){
				if(data.code==0){
					that.posdata.ordertime=(new Date()).format("yyyy-MM-dd hh:mm:ss");
					// that.Print(function(){
					// 	that.success("提交成功！");
					// 	$btn.button('reset');
					// 	that.refresh();
					//
					// });
                    $btn.button('reset');
                    that.success("提交成功！");
                    that.refresh();
				}else{
					that.error("提交失败！",data.msg)
					$btn.button('reset');
				}
			},'json');
		},

		isIE8:function(){
			return navigator.appName == "Microsoft Internet Explorer" && navigator.appVersion.match(/8./i)=="8."
		},

		refresh:function(){
 			if(this.isIE8()){
 				window.location.href=window.location.href;
 				return;
 			}
			this.posdata = {
				orderno :"",
				customno :"",
				totamt:0,
				trueamt:0,
				payamt:0,
				retamt:0,
				desc:"",
				items:new Array()
			};
			this.loadOrderno();
			this.tot_changed = true;
			this.item_changed = true;
			this.drawHtml();
		},
		
		// Print:function(func){
		// 	var data = {};
		// 	data.orderno = this.posdata.orderno;
		// 	var itemshtml = "";
		// 	for(var i =0;i<this.posdata.items.length;i++){
		// 		itemshtml = itemshtml+nano("<tr><td>{product_desc}</td><td>{qty}</td><td>{trueprice}</td><td>{amt}</td></tr>",this.posdata.items[i]);
		// 	}
		// 	data.itemhtml = itemshtml;
		// 	data.ordertime = this.posdata.ordertime;
		// 	data.totamt = this.posdata.totamt;
		// 	data.trueamt = this.posdata.trueamt;
		// 	data.payamt = this.posdata.payamt;
		// 	data.retamt = this.posdata.retamt;
		// 	$("#page1").html(nano($("#page1-template").html(),data));
		// 	doPrint();
		// 	if(this.isIE8()){
		// 		setTimeout(function(){
		// 			func&&func();
		// 		},3000)
 		// 	}else{
 		// 		func&&func();
 		// 	}
		//
		// },
		success:function(msg,msg2){
			toastr.success(msg,msg2);
		},
		error:function(msg,msg2){
			toastr.error(msg,msg2);
		}
	};

	var pos = window.pos = new Pos();
	
	$(document).ready(function(e){
		pos.init();

	});
})($);