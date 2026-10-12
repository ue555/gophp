<?php
echo (false && false) + 0, (false && true) + 0, (true && false) + 0, (true && true) + 0, '|';
echo (false || false) + 0, (false || true) + 0, (true || false) + 0, (true || true) + 0, '|';
echo (!true) + 0, (!false) + 0, "\n";
