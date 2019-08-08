// Code generated by "stringer -type=Type"; DO NOT EDIT.

package payload

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeUnknown-0]
	_ = x[TypeMeta-1]
	_ = x[TypeError-2]
	_ = x[TypeID-3]
	_ = x[TypeIDs-4]
	_ = x[TypeJet-5]
	_ = x[TypeState-6]
	_ = x[TypeGetObject-7]
	_ = x[TypePassState-8]
	_ = x[TypeIndex-9]
	_ = x[TypePass-10]
	_ = x[TypeGetCode-11]
	_ = x[TypeCode-12]
	_ = x[TypeSetCode-13]
	_ = x[TypeSetIncomingRequest-14]
	_ = x[TypeSetOutgoingRequest-15]
	_ = x[TypeSagaCallAcceptNotification-16]
	_ = x[TypeGetFilament-17]
	_ = x[TypeGetRequest-18]
	_ = x[TypeRequest-19]
	_ = x[TypeFilamentSegment-20]
	_ = x[TypeSetResult-21]
	_ = x[TypeActivate-22]
	_ = x[TypeRequestInfo-23]
	_ = x[TypeGotHotConfirmation-24]
	_ = x[TypeDeactivate-25]
	_ = x[TypeUpdate-26]
	_ = x[TypeHotObjects-27]
	_ = x[TypeResultInfo-28]
	_ = x[TypeGetPendings-29]
	_ = x[TypeHasPendings-30]
	_ = x[TypePendingsInfo-31]
	_ = x[TypeReplication-32]
	_ = x[TypeGetJet-33]
	_ = x[TypeAbandonedRequestsNotification-34]
	_ = x[TypeGetLightInitialState-35]
	_ = x[TypeLightInitialState-36]
	_ = x[TypeGetIndex-37]
	_ = x[TypeUpdateJet-38]
	_ = x[TypeReturnResults-39]
	_ = x[TypeCallMethod-40]
	_ = x[TypeExecutorResults-41]
	_ = x[TypePendingFinished-42]
	_ = x[TypeAdditionalCallFromPreviousExecutor-43]
	_ = x[TypeStillExecuting-44]
	_ = x[TypeOk-45]
	_ = x[_latestType-46]
}

const _Type_name = "TypeUnknownTypeMetaTypeErrorTypeIDTypeIDsTypeJetTypeStateTypeGetObjectTypePassStateTypeIndexTypePassTypeGetCodeTypeCodeTypeSetCodeTypeSetIncomingRequestTypeSetOutgoingRequestTypeSagaCallAcceptNotificationTypeGetFilamentTypeGetRequestTypeRequestTypeFilamentSegmentTypeSetResultTypeActivateTypeRequestInfoTypeGotHotConfirmationTypeDeactivateTypeUpdateTypeHotObjectsTypeResultInfoTypeGetPendingsTypeHasPendingsTypePendingsInfoTypeReplicationTypeGetJetTypeAbandonedRequestsNotificationTypeGetLightInitialStateTypeLightInitialStateTypeGetIndexTypeUpdateJetTypeReturnResultsTypeCallMethodTypeExecutorResultsTypePendingFinishedTypeAdditionalCallFromPreviousExecutorTypeStillExecutingTypeOk_latestType"

var _Type_index = [...]uint16{0, 11, 19, 28, 34, 41, 48, 57, 70, 83, 92, 100, 111, 119, 130, 152, 174, 204, 219, 233, 244, 263, 276, 288, 303, 325, 339, 349, 363, 377, 392, 407, 423, 438, 448, 481, 505, 526, 538, 551, 568, 582, 601, 620, 658, 676, 682, 693}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
