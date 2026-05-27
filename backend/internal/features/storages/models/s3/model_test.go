package s3_storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculatePartSize_WhenFileJustBelowMinThreshold_ReturnsMinPartSize(t *testing.T) {
	// partSize falls below MinPartSize iff fileSize < MinPartSize * maxParts.
	// MinPartSize * maxParts = 5 MiB * 1000 ≈ 4.88 GiB, so 4 GiB stays clamped.
	fourGB := int64(4 * 1024 * 1024 * 1024)

	result := CalculatePartSize(fourGB)

	assert.Equal(t, MinPartSize, result)
}

func Test_CalculatePartSize_WhenFileFifteenGB_ScalesUpChunks(t *testing.T) {
	fifteenGB := int64(15 * 1024 * 1024 * 1024)

	result := CalculatePartSize(fifteenGB)

	// ceil(15 GB / 1000) ≈ 16.1 MB
	expectedMin := fifteenGB / maxParts
	assert.Greater(t, result, expectedMin)
	assert.GreaterOrEqual(t, result, MinPartSize)
	assert.LessOrEqual(t, result, MaxPartSize)
	// Verify the calculated part size can cover the file in ≤ 1000 parts
	assert.LessOrEqual(t, (fifteenGB+result-1)/result, int64(maxParts))
}

func Test_CalculatePartSize_WhenFileSizeUnknown_ReturnsStreamingDefault(t *testing.T) {
	assert.Equal(t, DefaultPartSize, CalculatePartSize(0))
	assert.Equal(t, DefaultPartSize, CalculatePartSize(-1))
}

func Test_CalculatePartSize_WhenFileEmpty_ReturnsStreamingDefault(t *testing.T) {
	result := CalculatePartSize(0)

	assert.Equal(t, DefaultPartSize, result)
}

func Test_CalculatePartSize_WhenFileExactlyFiveGB_ScalesAboveMinPartSize(t *testing.T) {
	exactlyFiveGB := int64(5 * 1024 * 1024 * 1024)

	result := CalculatePartSize(exactlyFiveGB)

	// ceil(5 GiB / 1000) ≈ 5.12 MiB, just above MinPartSize (5 MiB),
	// so the function returns the computed chunk size rather than clamping.
	assert.Greater(t, result, MinPartSize)
	assert.LessOrEqual(t, result, MaxPartSize)
	assert.LessOrEqual(t, (exactlyFiveGB+result-1)/result, int64(maxParts))
}

func Test_CalculatePartSize_WhenFileSmall_ReturnsMinPartSize(t *testing.T) {
	oneMB := int64(1 * 1024 * 1024)

	result := CalculatePartSize(oneMB)

	assert.Equal(t, MinPartSize, result)
}

func Test_CalculatePartSize_WhenFileNearMaxSize_RespectsMaxPartSize(t *testing.T) {
	// 5 TB file — near the absolute max of S3 (5 TB object limit)
	fiveTB := int64(5 * 1024 * 1024 * 1024 * 1024)

	result := CalculatePartSize(fiveTB)

	assert.Equal(t, MaxPartSize, result)
}

func Test_getPartSize_WhenS3PartSizeConfigured_UsesConfiguredValue(t *testing.T) {
	s := &S3Storage{S3PartSize: 64 * 1024 * 1024} // 64 MB

	result := s.getPartSize()

	assert.Equal(t, int64(64*1024*1024), result)
}

func Test_getPartSize_WhenS3PartSizeZero_ReturnsStreamingDefault(t *testing.T) {
	s := &S3Storage{S3PartSize: 0}

	result := s.getPartSize()

	assert.Equal(t, DefaultPartSize, result)
}

func Test_getPartSize_WhenS3PartSizeBelowMinimum_ClampsToMinimum(t *testing.T) {
	s := &S3Storage{S3PartSize: 1024} // 1 KB — way below minimum

	result := s.getPartSize()

	assert.Equal(t, MinPartSize, result)
}

func Test_getPartSize_WhenS3PartSizeAboveMaximum_ClampsToMaximum(t *testing.T) {
	s := &S3Storage{S3PartSize: 10 * 1024 * 1024 * 1024} // 10 GB — above maximum

	result := s.getPartSize()

	assert.Equal(t, MaxPartSize, result)
}
