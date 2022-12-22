package events

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
	eventSerialize "ws/framework/plugin/event_serialize"
	functionTools "ws/framework/utils/function_tools"
)

var moduleHexInit sync.Once
var moduleHexStr string

func moduleHex() string {
	moduleHexInit.Do(func() {
		// libsubstitute substitute-inserter substitute-loader libsubstrate
		// libWhatsAppDylib libcycript RevealServer Dobby
		imageList := "WhatsApp,libc++.1.dylib,LinkPresentation,libSystem.B.dylib,libcompression.dylib,CoreBluetooth,Core,Network,libz.1.dylib,MetalKit,libsqlite3.dylib,libxml2.2.dylib,MetricKit,EventKit,PushKit,BackgroundTasks,EventKitUI,SafariServices,VideoToolbox,Foundation,libobjc.A.dylib,AVFoundation,AudioToolbox,CallKit,CloudKit,Contacts,ContactsUI,CoreData,CoreFoundation,CoreGraphics,CoreImage,CoreLocation,CoreMedia,CoreMotion,CoreServices,CoreSpotlight,CoreText,CoreVideo,GLKit,ImageIO,Intents,LocalAuthentication,MapKit,MediaPlayer,MessageUI,Metal,OpenGLES,PassKit,Photos,PhotosUI,QuartzCore,Security,StoreKit,UIKit,UserNotifications,Vision,WebKit,libswiftCoreMIDI.dylib,libswiftCoreML.dylib,libswiftAVFoundation.dylib,libswiftAccelerate.dylib,libswiftAssetsLibrary.dylib,libswiftCallKit.dylib,libswiftCloudKit.dylib,libswiftContacts.dylib,libswiftCore.dylib,libswiftCoreAudio.dylib,libswiftCoreData.dylib,libswiftCoreFoundation.dylib,libswiftCoreGraphics.dylib,libswiftCoreImage.dylib,libswiftCoreLocation.dylib,libswiftCoreMedia.dylib,libswiftDarwin.dylib,libswiftDispatch.dylib,libswiftFoundation.dylib,libswiftGLKit.dylib,libswiftIntents.dylib,libswiftMapKit.dylib,libswiftMediaPlayer.dylib,libswiftMetal.dylib,libswiftModelIO.dylib,libswiftNetwork.dylib,libswiftObjectiveC.dylib,libswiftPhotos.dylib,libswiftQuartzCore.dylib,libswiftUIKit.dylib,libswiftVision.dylib,libswift_Concurrency.dylib,libswiftos.dylib,libswiftsimd.dylib,libc++abi.dylib,libcache.dylib,libcommonCrypto.dylib,libcompiler_rt.dylib,libcopyfile.dylib,libcorecrypto.dylib,libdispatch.dylib,libdyld.dylib,liblaunch.dylib,libmacho.dylib,libremovefile.dylib,libsystem_asl.dylib,libsystem_blocks.dylib,libsystem_c.dylib,libsystem_configuration.dylib,libsystem_containermanager.dylib,libsystem_coreservices.dylib,libsystem_darwin.dylib,libsystem_dnssd.dylib,libsystem_featureflags.dylib,libsystem_info.dylib,libsystem_m.dylib,libsystem_malloc.dylib,libsystem_networkextension.dylib,libsystem_notify.dylib,libsystem_sandbox.dylib,libsystem_kernel.dylib,libsystem_platform.dylib,libsystem_pthread.dylib,libsystem_symptoms.dylib,libsystem_trace.dylib,libunwind.dylib,libxpc.dylib,URLFormatting,IconServices,Accelerate,AggregateDictionary,AppleMediaServices,AssertionServices,AVKit,Celestial,MediaRemote,QuickLookThumbnailing,TCC,libMobileGestalt.dylib,CFNetwork,libicucore.A.dylib,IOKit,SystemConfiguration,libnetwork.dylib,libapple_nghttp2.dylib,IOMobileFramebuffer,libbsm.0.dylib,libpcap.A.dylib,libcoretls.dylib,libcoretls_cfhelpers.dylib,libenergytrace.dylib,libarchive.2.dylib,libCRFSuite.dylib,liblangid.dylib,liblzma.5.dylib,libbz2.1.0.dylib,libiconv.2.dylib,libcharset.1.dylib,IOSurface,AddressBookLegacy,AppSupport,ManagedConfiguration,MobileCoreServices,libCTGreenTeaLogger.dylib,CellularPlanManager,ClassKit,CoreSuggestions,PhoneNumbers,vCard,ContactsFoundation,Accounts,DataAccessExpress,CorePhoneNumbers,BaseBoard,RunningBoardServices,PersistentConnection,ProtocolBuffer,CoreTelephony,CommonUtilities,libcupolicy.dylib,libTelephonyUtilDynamic.dylib,MobileInstallation,CoreServicesStore,MobileSystemServices,MobileWiFi,Bom,MobileKeyBag,CaptiveNetwork,EAP8021X,CoreAnalytics,APFS,AppleSauce,libutil.dylib,libFontParser.dylib,libhvf.dylib,vImage,vecLib,libvMisc.dylib,libvDSP.dylib,libBLAS.dylib,libLAPACK.dylib,libLinearAlgebra.dylib,libSparseBLAS.dylib,libQuadrature.dylib,libBNNS.dylib,libSparse.dylib,IOSurfaceAccelerator,libate.dylib,AppleJPEG,IOAccelerator,libCoreFSCache.dylib,liblockdown.dylib,libmis.dylib,Netrb,GraphicsServices,DataMigration,UserManagement,libGSFont.dylib,FontServices,libGSFontCache.dylib,libAccessibility.dylib,OTSVG,ConstantClasses,AXCoreUtilities,MediaAccessibility,SpringBoardServices,IdleTimerServices,BoardServices,FrontBoardServices,BackBoardServices,DataDetectorsCore,CoreNLP,AppleFSCompression,libmecab.dylib,libgermantok.dylib,libThaiTokenizer.dylib,libChineseTokenizer.dylib,LanguageModeling,CoreEmoji,LinguisticData,Lexicon,libcmph.dylib,CloudDocs,AppleAccount,CrashReporterSupport,SymptomDiagnosticReporter,ApplePushService,C2,ProtectedCloudStorage,CoreUtils,OSAnalytics,CoreSymbolication,Symbolication,OSAServicesClient,MallocStackLogging,AuthKit,NetworkExtension,DeviceIdentity,AppleIDAuthSupport,SharedUtils,libnetworkextension.dylib,GeoServices,LocationSupport,CoreLocationProtobuf,PowerLog,NanoPreferencesSync,NanoRegistry,libheimdal-asn1.dylib,FileProvider,MobileSpotlightIndex,GenerationalStorage,ChunkingLibrary,MetadataUtilities,libprequelite.dylib,MobileIcons,CoreUI,libFosl_dynamic.dylib,ColorSync,GraphVisualizer,MetalPerformanceShaders,FaceCore,libncurses.5.4.dylib,WatchdogClient,CoreAudio,ASEProcessing,libtailspin.dylib,libEDR,SignpostCollection,ktrace,SampleAnalysis,kperfdata,libdscsym.dylib,SignpostSupport,LoggingSupport,kperf,CoreBrightness,libIOReport.dylib,CPMS,HID,libGFXShared.dylib,libGLImage.dylib,libCVMSPluginSupport.dylib,libCoreVMClient.dylib,MPSCore,MPSImage,MPSNeuralNetwork,MPSMatrix,MPSRayIntersector,MPSNDArray,AudioToolboxCore,caulk,libAudioToolboxUtility.dylib,MediaExperience,TextureIO,CoreSVG,InternationalSupport,SetupAssistant,AppleIDSSOAuthentication,AccountSettings,CoreFollowUp,SetupAssistantSupport,MobileBackup,CoreTime,IntlPreferences,AppConduit,Rapport,StreamingZip,MobileDeviceLink,AccountsDaemon,GSS,PlugInKit,IDS,WirelessDiagnostics,OAuth,Heimdal,libresolv.9.dylib,CommonAuth,Marco,IMFoundation,IDSFoundation,Engram,libtidy.A.dylib,libAWDSupportFramework.dylib,libAWDSupport.dylib,libprotobuf-lite.dylib,libprotobuf.dylib,ProactiveEventTracker,ProactiveSupport,SearchFoundation,DataDetectorsNaturalLanguage,IntentsFoundation,InternationalTextSearch,MobileAsset,ResponseKit,CalendarDaemon,CalendarDatabase,CalendarFoundation,iCalendar,BackgroundTaskAgent,PersonaKit,CoreDAV,NLP,Montreal,CryptoTokenKit,AVFAudio,MediaToolbox,Quagga,libAudioStatistics.dylib,perfdata,libperfcheck.dylib,CoreAUC,CoreHaptics,NetworkStatistics,Pegasus,UIKitServices,DocumentManager,UIKitCore,ShareSheet,DocumentManagerCore,UIFoundation,XCTTargetBootstrap,WebKitLegacy,SAObjects,HangTracer,SignpostMetrics,PointerUIServices,StudyLog,CoreMaterial,libapp_launch_measurement.dylib,PhysicsKit,PrototypeTools,TextInput,JavaScriptCore,WebCore,libwebrtc.dylib,RemoteTextInput,MediaServices,CorePrediction,PDFKit,SafariSafeBrowsing,CoreOptimization,CorePDF,RevealCore,NaturalLanguage,CoreML,DuetActivityScheduler,Espresso,CoreDuet,CoreDuetContext,CoreDuetDebugLogging,CoreDuetDaemonProtocol,AppleNeuralEngine,ANECompiler,ANEServices,libsandbox.1.dylib,libMatch.1.dylib,CTCarrierSpace,IntentsUI,QuickLook,PassKitCore,PassKitUI,AddressBook,BridgePreferences,CoreRecognition,PBBridgeSupport,HSAAuthentication,CoreRecents,EmailCore,DigitalAccess,CoreIDV,DifferentialPrivacy,RTCReporting,CoreCDP,EmailFoundation,SEService,CloudServices,KeychainCircle,Messages,SpringBoardUIServices,AppleAccountUI,OnBoardingKit,RemoteUI,LocalAuthenticationPrivateUI,CoreCDPUI,StoreServices,AuthKitUI,PassKitUIFoundation,Preferences,FindMyDevice,BluetoothManager,Sharing,ChatKit,IMSharedUtilities,PhotosPlayer,EmojiKit,PhotosFormats,ToneLibrary,CameraKit,StorageSettings,ScreenTimeCore,ContactsAutocompleteUI,IMTranscoding,TelephonyUtilities,IMSharedUI,IMDPersistence,IMCore,FTServices,FTClientServices,SettingsFoundation,MobileBluetooth,BiometricKit,ExposureNotification,FamilyCircle,AccountsUI,Message,ManagedConfigurationUI,CertInfo,CertUI,PhotoLibraryServices,NanoResourceGrabber,NanoSystemSettings,MediaConversionService,DCIMServices,MediaStream,PhotoFoundation,CloudPhotoLibrary,CacheDelete,CloudPhotoServices,CoreMediaStream,AssetsLibraryServices,PhotosImagingFoundation,MMCS,AssetCacheServices,FSEvents,DiagnosticLogCollection,WatchKit,DoNotDisturb,IDSKVStore,DeviceManagement,CoreSDB,CommunicationsFilter,IncomingCallFilter,libsysdiagnose.dylib,FTAWD,MobileActivation,Catalyst,Categories,ContextKit,AppSupportUI,PersonaUI,ContactsDonation,ContactsUICore,CoreCDPInternal,FMCoreLite,MailSupport,DataAccess,Email,EmailAddressing,EmailDaemon,MailServices,MessageSupport,MIME,Notes,EmojiFoundation,ImageCaptureCore,QuickLookSupport,HomeSharing,iTunesCloud,MediaLibraryCore,MediaPlatform,MusicLibrary,DAAPKit,ConfigurationEngineModel,ContactsAutocomplete,VisionKit,PencilKit,Futhark,TextRecognition,Navigation,VectorKit,CoreHandwriting,MaterialKit,UIAccessibility,ProofReader,libmecabra.dylib,BaseBoardUI,AXMediaUtilities,libAXSpeechManager.dylib,AccessibilityUIUtilities,ScreenReaderCore,libAXSafeCategoryBundle.dylib,AccessibilityUtilities,AXRuntime,TextToSpeech,CoreMIDI,AccessibilitySharedSupport,AssistantServices,SiriInstrumentation,VoiceServices,SiriTTS,libedit.3.dylib,AppStoreDaemon,SpringBoardFoundation,TelephonyUI,ProgressUI,MobileObliteration,PrototypeToolsUI,CalendarUIKit,SceneKit,ModelIO,SharedWebCredentials,AuthenticationServices,WebBookmarks,SafariSharedUI,SafariFoundation,WebUI,SafariCore,AppSSO,PersonalizationPortrait,UserActivity,CoreParsec,ParsecSubscriptionServiceSupport,SafariShared,VideoSubscriberAccount,MetricsKit,iTunesStore,Social,UsageTracking,AssetsLibrary,IntentsCore,Combine,PhotoLibrary,PhotosUICore,CameraEditKit,NeutrinoCore,NeutrinoKit,PhotoImaging,CPAnalytics,VideoProcessing,AppleCVA,AutoLoop,CoreAppleCVA,SoundAnalysis,InertiaCam,DistributedEvaluation,CloudKitCodeProtobuf,CryptoKit,CloudKitCode,CryptoKitCBridging,CryptoKitPrivate,libobjc-trampolines.dylib"

		uuids := strings.Split(imageList, ",")

		s256 := sha256.New()
		buffer := bytes.Buffer{}
		hexEncoder := hex.NewEncoder(&buffer)

		for _, uuid := range uuids {
			s256.Reset()

			s256.Write(functionTools.S2B(uuid))

			_, _ = hexEncoder.Write(s256.Sum(nil))
		}

		s256.Reset()
		s256.Write(buffer.Bytes())

		moduleHexStr = hex.EncodeToString(s256.Sum(nil))
	})

	return moduleHexStr
}

type WamEventDaily struct {
	WAMessageEvent

	IphoneNotificationAuthStatus float64
	OSNotifySetting              float64
	AccessVoiceover              float64
	AddressBookSize              float64
	AddressBookWSASize           float64
	AudioCellular                float64
	AudioWIFI                    float64
	DocCellular                  float64
	DocWifi                      float64
	ImageCellular                float64
	ImageWIFI                    float64
	VideoCellular                float64
	VideoWifi                    float64
	BackupRestoreEnc             float64
	ChatDBSize                   float64
	DBSearchFTS                  float64

	FavoriteStickerCount      float64
	FavoriteFirstStickerCount float64
	FavoriteTotalStickerCount float64
	InstallStickerCount       float64
	Install3rdStickerCount    float64
	InstallFirstStickerCount  float64
	InstallTotalStickerCount  float64

	ICloudBackupInterval float64
	ICloudSignedIn       float64
	JailBroken           float64
	ModuleHash           string
	AlertEnabled         float64
	BadgesEnabled        float64
	SoundsEnabled        float64
	Language             string
	Location             string
	PackageName          string
	PaymentEnabled       float64
	StorageSize          float64
	StorageCacheSize     float64
	StorageTotalSize     float64
}

type DailyOption struct {
	Language string
	Location string
}

func WithDailyOption(lang, location string) DailyOption {
	return DailyOption{
		Language: lang,
		Location: location,
	}
}

func (event *WamEventDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x82, event.IphoneNotificationAuthStatus).
		SerializeNumber(0x76, event.OSNotifySetting).
		SerializeNumber(0x6c, event.AccessVoiceover).
		SerializeNumber(0xb, event.AddressBookSize).
		SerializeNumber(0xc, event.AddressBookWSASize).
		SerializeNumber(0x5a, event.AudioCellular).
		SerializeNumber(0x59, event.AudioWIFI).
		SerializeNumber(0x60, event.DocCellular).
		SerializeNumber(0x5f, event.DocWifi).
		SerializeNumber(0x57, event.ImageCellular).
		SerializeNumber(0x56, event.ImageWIFI).
		SerializeNumber(0x5d, event.VideoCellular).
		SerializeNumber(0x5c, event.VideoWifi).
		SerializeNumber(0x8a, event.BackupRestoreEnc).
		SerializeNumber(0x13, event.ChatDBSize).
		SerializeNumber(0x1e, event.DBSearchFTS).
		SerializeNumber(0x71, event.FavoriteStickerCount).
		SerializeNumber(0x70, event.FavoriteFirstStickerCount).
		SerializeNumber(0x6f, event.FavoriteTotalStickerCount).
		SerializeNumber(0x74, event.InstallStickerCount).
		SerializeNumber(0x89, event.Install3rdStickerCount).
		SerializeNumber(0x73, event.InstallFirstStickerCount).
		SerializeNumber(0x72, event.InstallTotalStickerCount).
		SerializeNumber(0x1, event.ICloudBackupInterval).
		SerializeNumber(0x1d, event.ICloudSignedIn).
		SerializeNumber(0x1c, event.JailBroken).
		SerializeString(0x75, event.ModuleHash).
		SerializeNumber(0x83, event.AlertEnabled).
		SerializeNumber(0x85, event.BadgesEnabled).
		SerializeNumber(0x84, event.SoundsEnabled).
		SerializeString(0x5, event.Language).
		SerializeString(0x6, event.Location).
		SerializeString(0x66, event.PackageName).
		SerializeNumber(0x64, event.PaymentEnabled).
		SerializeNumber(0x1f, event.StorageSize).
		SerializeNumber(0x88, event.StorageCacheSize)

	buffer.Footer().
		SerializeNumber(0x20, event.StorageTotalSize)
}

func (event *WamEventDaily) InitFields(option interface{}) {
	event.IphoneNotificationAuthStatus = 1
	event.OSNotifySetting = 1
	event.AccessVoiceover = 0
	event.AddressBookSize = 0
	event.AddressBookWSASize = 0
	event.AudioCellular = 0
	event.AudioWIFI = 1
	event.DocCellular = 0
	event.DocWifi = 1
	event.ImageCellular = 1
	event.ImageWIFI = 1
	event.VideoCellular = 0
	event.VideoWifi = 1
	event.BackupRestoreEnc = 6
	event.ChatDBSize = float64(rand.Intn(300000) + 300000)
	event.DBSearchFTS = 0

	event.FavoriteStickerCount = 0
	event.FavoriteFirstStickerCount = 0
	event.FavoriteTotalStickerCount = 0
	event.InstallStickerCount = 0
	event.Install3rdStickerCount = 0
	event.InstallFirstStickerCount = 0
	event.InstallTotalStickerCount = 0

	event.ICloudBackupInterval = -1
	event.ICloudSignedIn = 0
	event.JailBroken = 0 // 1:越狱 0:非越狱
	event.ModuleHash = moduleHex()
	event.AlertEnabled = 1
	event.BadgesEnabled = 1
	event.SoundsEnabled = 1
	//event.Language = "zh"
	//event.Location = "CN"
	event.PackageName = "net.whatsapp.WhatsApp"
	event.PaymentEnabled = 0
	event.StorageSize = 8338581430
	event.StorageCacheSize = 6873706496
	event.StorageTotalSize = 15989485568

	if opt, ok := option.(DailyOption); ok {
		event.Language = opt.Language
		event.Location = opt.Location
	}
}
