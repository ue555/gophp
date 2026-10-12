<?php
$equal = 1===1;
$different = 1!=='1';
echo ($equal === true) + 0, ($different === true) + 0, ($equal === 1) + 0, "\n";
echo ((1 !== 1) === false) + 0, ((null === false) === false) + 0, "\n";
