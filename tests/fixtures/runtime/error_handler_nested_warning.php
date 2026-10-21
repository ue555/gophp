<?php
function handlerWithWarning($level) {
    echo $inside;
    echo "handled|";
    return true;
}
set_error_handler('handlerWithWarning');
echo $outside;
echo "done\n";
