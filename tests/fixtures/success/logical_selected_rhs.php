<?php
echo (true && '0') === false;
echo (false || 'hello') === true;
echo (true && 0) === false;
echo (false || 1) === true;
echo (false || (true && !false)) === true;
echo "\n";
