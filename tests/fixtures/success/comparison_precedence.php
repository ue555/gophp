<?php
echo (1 + 2 < 4 * 2) === true;
echo (1 . 0 < 2 . 0) === true;
echo (1 < 2 == true) === true;
echo (false == 2 < 1) === true;
echo (1 < 2 === true) === true;
echo (1 >= 2 !== true) === true;
echo (1 < 2 == 3 >= 3) === true;
echo (!1 < 2) === true;
echo (1 == 1 && 2 != 3 || false) === true;
echo (1 > 2 || 3 <= 3 && 0 == "0") === true;
echo ((1 < 2) < 3) === false;
echo (1 < (2 < 3)) === false;
echo ((1 == 1) != false) === true;
echo (1 == (1 != 2)) === true;
echo (1 == "1") !== (1 === "1");
echo (!(1 == "2") && 1 !== "1") === true;
echo (false || 0 != "word" && 1 <= 1) === true;
echo "\n";
