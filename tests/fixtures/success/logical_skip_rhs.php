<?php
echo (false && $and_missing) === false;
echo (true || $or_missing) === true;
echo (0 && (1 / 0)) === false;
echo ('word' || (1 / 0)) === true;
echo ('' && ('bad' + 1)) === false;
echo (true || ('bad' + 1)) === true;
echo (false && !$not_missing) === false;
echo (true || ($nested_missing && (1 / 0))) === true;
echo "\n";
