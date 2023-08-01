(function () {
    var stopTimerTag = false;

    function stopTimer() {
        stopTimerTag = true;
    }

    window.getHeight=function () {
        return document.body.clientHeight;
    };

    function loaded() {
        setTimeout(function() {
            if (window.parent) window.parent.setMailContentHeight(window.name);
            if (!stopTimerTag) {
                loaded();
            } else {
                window.parent.scroll(0, 0);
            }
        }, 100);
    }

    function getTop(e){
        var offset = e.offsetTop;
        if(e.offsetParent != null) offset+=getTop(e.offsetParent);
        return offset;
    }

    function scrollDir(e) {
        var target = e.currentTarget;
        var id = target.href.split('#')[1] || '';
        var top = 0;
        var el = null;
        if (id) {
            el = document.getElementById(id);
            top = getTop(el) + 50;
        }
        window.parent.scroll(0, top);
        return false;
    }

    window.onload = stopTimer;
    loaded();
    if (window.navigator.userAgent.indexOf("Firefox") > 0) {
        var dirs = document.querySelectorAll('#dir-list a');
        dirs = Array.prototype.slice.call(dirs);
        if (dirs.length > 0) {
            for (var k = 0, dLen = dirs.length; k < dLen; k++) {
                dirs[k].addEventListener('click', scrollDir, false);
            }
        }
    }
    jQuery('#switch-dir').on('click', function () {
        var self = $(this),
            dirContent = $('#dir-content');
        if (dirContent.hasClass('hide')) {
            self.find('.text').text('隐藏');
        } else {
            self.find('.text').text('显示');
        }
        dirContent.toggleClass('hide');
    });
    // iframe中的a标签打开链接
    // $("a").click(function(){
	// 	var url = $(this).attr('href');
	// 	if(url.indexOf('/html/api_detail_')>-1){
	// 		parent.redirect(url);
	// 	}
    // });
})();