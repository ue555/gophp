<?php
echo ((true || false) && false) === false;
echo (false || false || '0' || 'yes') === true;
echo (true && 'yes' && 1) === true;
echo (true && '0' && $missing) === false;
echo (!!'0') === false;
echo (!!!'0') === true;
echo (!true!==false) === false;
echo (!(1 === 1 && (false || 2))) === false;
echo "\n";
