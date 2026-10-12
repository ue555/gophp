<?php
echo (true || false && false) === true;
echo (false || true && false) === false;
echo (!1 === false) === true;
echo (!0 + 2 === 3) === true;
echo (!2 * 3 === 0) === true;
echo (!0 . 'x' === '1x') === true;
echo ('x' . false && true) === true;
echo ('x' . false || false) === true;
echo (1 + 2 === 3 && 0 . '0' === '00') === true;
echo (1 !== 1 || 2 * 3 === 6) === true;
echo (!-1) === false;
echo (-!false === -1) === true;
echo (!+0) === true;
echo (false === false && true === true) === true;
echo "\n";
