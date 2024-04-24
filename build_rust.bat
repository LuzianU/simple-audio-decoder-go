cd simple-audio-decoder-rs
git pull
cargo build --release
cd ..
xcopy /Y simple-audio-decoder-rs\target\release\simple_audio_decoder_rs.dll .\