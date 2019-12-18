document.addEventListener('DOMContentLoaded', function() {

    $(".show-comment-form").click(function() {
        $("#wrap").slideToggle();
    });
    var aluContainer = document.querySelector('.comment-form-smilies');

    $('#tit span').click(function() {
        console.log("click", this, "$", $(this), "下标", $(this).index())
        var i = $(this).index(); //下标第一种写法
        //var i = $('tit').index(this);//下标第二种写法
        $(this).addClass('select').siblings().removeClass('select');
        console.log("show", $('#con li'))
        $('#con li').eq(i).show().siblings().hide();
        aluContainer = document.querySelector('.comment-form-smilies');
        console.log("判断", aluContainer)
    });
    // $('.add-smily').click(function() {
    //     console.log("smile", this);

    // });

    if (!aluContainer) return;
    $('.add-smily').click(function(e) {
        console.log("click", e)
        var myField,
            _self = e.target.dataset.smilies ? e.target : e.target.parentNode;
        if (typeof _self.dataset.smilies == 'undefined') return;
        var tag = ' ' + _self.dataset.smilies + ' ';
        if (document.getElementById('comment') && document.getElementById('comment').type == 'textarea') {
            myField = document.getElementById('comment')
        } else {
            return false
        }
        if (document.selection) {
            myField.focus();
            sel = document.selection.createRange();
            sel.text = tag;
            myField.focus()
        } else if (myField.selectionStart || myField.selectionStart == '0') {
            var startPos = myField.selectionStart;
            var endPos = myField.selectionEnd;
            var cursorPos = endPos;
            myField.value = myField.value.substring(0, startPos) + tag + myField.value.substring(endPos, myField.value.length);
            cursorPos += tag.length;
            myField.focus();
            myField.selectionStart = cursorPos;
            myField.selectionEnd = cursorPos
        } else {
            myField.value += tag;
            myField.focus()
        }
    });
});