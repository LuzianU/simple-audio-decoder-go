git submodule update --init --force --remote
cd simple-audio-decoder-rs
cargo build --release
cd ..
xcopy /Y simple-audio-decoder-rs\target\release\simple_audio_decoder_rs.dll .\
xcopy /Y simple-audio-decoder-rs\src\simple_audio_decoder_rs.h .\
