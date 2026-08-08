package rtp

import (
	"testing"
	"time"

	"github.com/langhuihui/gomem"
	"github.com/stretchr/testify/require"
	. "m7s.live/v5/pkg"
	"m7s.live/v5/pkg/codec"
)

func TestG711MuxSequenceContinuesAcrossFrames(t *testing.T) {
	tests := []struct {
		name        string
		codec       IAudioCodecCtx
		payloadType uint8
	}{
		{name: "PCMA", codec: codec.NewPCMACtx(), payloadType: 8},
		{name: "PCMU", codec: codec.NewPCMUCtx(), payloadType: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allocator := gomem.NewScalableMemoryAllocator(4096)
			var output AudioFrame
			output.SetAllocator(allocator)

			for frameNumber := 0; frameNumber < 3; frameNumber++ {
				var audio AudioData
				audio.PushOne(make([]byte, 160))
				input := Sample{
					ICodecCtx: tt.codec,
					BaseSample: &BaseSample{
						Raw:       &audio,
						Timestamp: time.Duration(frameNumber) * 20 * time.Millisecond,
					},
				}

				require.NoError(t, output.Mux(&input))
				require.Len(t, output.Packets, 1)
				packet := output.Packets[0]
				require.Equal(t, uint16(frameNumber+1), packet.SequenceNumber)
				require.Equal(t, tt.payloadType, packet.PayloadType)
				require.Equal(t, uint32(frameNumber*160), packet.Timestamp)
				require.True(t, packet.Marker)

				output.Recycle()
			}
		})
	}
}
