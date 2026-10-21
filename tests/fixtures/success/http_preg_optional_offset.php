<?php
preg_match_all('/a(b)?(c)/', 'ac abc', $matches, PREG_SET_ORDER | PREG_OFFSET_CAPTURE);
echo json_encode($matches);
