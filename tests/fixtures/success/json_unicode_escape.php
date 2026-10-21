<?php
echo json_encode(['日本'=>'😀/']),"\n";
echo json_encode(['日本'=>'😀/'],JSON_UNESCAPED_UNICODE),"\n";
echo json_encode('日本😀/',JSON_UNESCAPED_UNICODE|JSON_UNESCAPED_SLASHES),"\n";
