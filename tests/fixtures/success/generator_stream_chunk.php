<?php
$stream = fopen('php://output', 'w');
stream_set_chunk_size($stream, 4096);
echo stream_set_chunk_size($stream, 1024), '|';
echo stream_set_chunk_size($stream, 2048), '|';
fclose($stream);
